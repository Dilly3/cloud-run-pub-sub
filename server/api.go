package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dilly3/cloud-run-pub-sub/config"
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
	//logger     *slog.Logger
}

type Context struct {
	*gin.Context
}

func handlerFunc(f func(*Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(&Context{c})
	}
}

func NewServer() *Server {
	config, err := config.GetConfig()
	if err != nil {
		log.Printf("Failed to get config: %v", err)
	}
	return &Server{Config: *config}
}

func (s *Server) SetupServer() {
	router := s.SetupRouter()

	server := &http.Server{
		Addr:    ":" + s.Config.Port,
		Handler: router,
	}

	s.httpServer = server
}

func (s *Server) SetupRouter() *gin.Engine {

	if s.Config.IsProd() {
		os.Setenv("GIN_MODE", "release")
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

	return router
}
