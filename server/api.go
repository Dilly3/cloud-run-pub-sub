package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dilly3/cloud-run-pub-sub/config"
	"github.com/dilly3/cloud-run-pub-sub/publisher"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey = "Authorization"
	RefreshTokenHeaderKey  = "X-Refresh-Token"
)

type Server struct {
	httpServer *http.Server
	Config     config.Configuration
	logger     *slog.Logger
	publisher  *publisher.Publisher
}

type Context struct {
	*gin.Context
}

func handlerFunc(f func(*Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(&Context{c})
	}
}

func NewServer(logger *slog.Logger, publisher *publisher.Publisher) *Server {
	config, err := config.GetConfig(logger)
	if err != nil {
		log.Printf("Failed to get config: %v", err)
	}
	return &Server{Config: *config, logger: logger, publisher: publisher}
}

func (s *Server) SetupServer(port string) *http.Server {
	router := s.SetupRouter()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	s.httpServer = server

	return server
}

func (s *Server) SetupRouter() *gin.Engine {

	if s.Config.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// add cors config
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // TODO: change this to the actual domain
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", AuthorizationHeaderKey, RefreshTokenHeaderKey, "Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		ExposeHeaders:    []string{"Content-Length", AuthorizationHeaderKey, "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiRouter := router.Group("/api/v1/")

	apiRouter.GET("/health", handlerFunc(s.HealthCheck))
	apiRouter.GET("/transactions", handlerFunc(s.GetTransactions))
	apiRouter.GET("/transactions/:id", handlerFunc(s.GetTransaction))
	apiRouter.GET("/publish", handlerFunc(s.PublishTransaction))
	apiRouter.POST("/transactions/poll", handlerFunc(s.PollTransaction))

	return router
}

func GracefulShutdown(srv *http.Server, logger *slog.Logger) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	logger.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "server error", err)
		os.Exit(1)
	}

	logger.Info("Server exiting")
}
