package server

import (
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/dilly3/cloud-run-pub-sub/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) QueueTransaction(c *Context) {
	// random transaction
	transaction := Transactions[rand.Intn(len(Transactions))]

	name, err := s.publisher.QueueTask(transaction, s.Config.TaskDelay)
	if err != nil {
		s.logger.Error("Failed to queue transaction", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction queued", "name": name})
}

func (s *Server) PublishTransaction(c *Context) {

	var transaction models.Transaction
	if err := s.publisher.DecodeQueueData(c.Request.Body, &transaction); err != nil {
		s.logger.Error("Failed to decode transaction", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode transaction"})
		return
	}

	id, err := s.publisher.Publish(transaction)
	if err != nil {
		s.logger.Error("Failed to publish transaction", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction published", "id": id})
}

func (s *Server) PollTransaction(c *Context) {
	// Read the request body from the push notification
	var transaction models.Transaction
	if err := s.publisher.DecodePubSubData(c.Request.Body, &transaction); err != nil {
		s.logger.Error("Failed to decode transaction", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode transaction"})
		return
	}
	if transaction.ID == 2 {
		s.logger.Error("Transaction ID is 2")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is 2, this is a test for failure"})
		return
	}

	// Log the transaction details
	child := s.logger.With(
		slog.Group("transaction_info",
			slog.Int("id", transaction.ID),
			slog.Int64("amount", transaction.Amount),
			slog.String("status", transaction.Status),
			slog.String("direction", transaction.Direction),
			slog.String("accountName", transaction.AccountName),
		),
	)
	child.Info("Transaction polled directly")

	c.JSON(http.StatusOK, gin.H{"message": "Transaction polled", "transaction": transaction})
}
