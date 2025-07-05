package server

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/dilly3/cloud-run-pub-sub/models"
	"github.com/gin-gonic/gin"
)

// decodePubSubData decodes base64 data from Pub/Sub message
func decodePubSubData(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func (s *Server) PublishTransaction(c *Context) {
	transaction := Transactions[rand.Intn(len(Transactions))]
	delayInSeconds, err := strconv.ParseInt(c.Query("delay"), 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse delay", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse delay"})
		return
	}

	// publish transaction with delay
	id, err := s.publisher.Publish(transaction, delayInSeconds)
	if err != nil {
		s.logger.Error("Failed to publish transaction", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction published", "id": id})
}

func (s *Server) PollTransaction(c *Context) {
	// Read the request body from the push notification
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.logger.Error("Failed to read request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// parse as a Pub/Sub push message
	var pubSubMessage models.PubSubMessage
	if err := json.Unmarshal(body, &pubSubMessage); err == nil {
		// This is a Pub/Sub push message, extract the data
		if pubSubMessage.Message.Data == "" {
			s.logger.Error("Pub/Sub message has no data")
			c.JSON(http.StatusBadRequest, gin.H{"error": "No data in message"})
			return
		}

		// Decode the base64 data from Pub/Sub message
		decodedData, err := decodePubSubData(pubSubMessage.Message.Data)
		if err != nil {
			s.logger.Error("Failed to decode base64 data from Pub/Sub message", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 data"})
			return
		}

		var transaction models.Transaction
		if err := json.Unmarshal(decodedData, &transaction); err != nil {
			s.logger.Error("Failed to unmarshal transaction from Pub/Sub data", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction data"})
			return
		}

		s.logger.Info("Transaction received from Pub/Sub",
			"messageId", pubSubMessage.Message.MessageID,
			"subscription", pubSubMessage.Subscription,
			"transactionId", transaction.ID,
			"amount", transaction.Amount,
			"status", transaction.Status,
		)

		// Return 200 to acknowledge the message
		c.JSON(http.StatusOK, gin.H{"message": "Transaction processed successfully"})
		return
	}

	// If not a Pub/Sub message, try to parse as direct transaction JSON
	var transaction models.Transaction
	if err := json.Unmarshal(body, &transaction); err != nil {
		s.logger.Error("Failed to unmarshal transaction", "error", err, "body", string(body))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
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

	// Process the transaction here
	// TODO: Add your business logic

	c.JSON(http.StatusOK, gin.H{"message": "Transaction polled", "transaction": transaction})
}
