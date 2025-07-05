package server

import (
	"net/http"
	"strconv"

	"github.com/dilly3/cloud-run-pub-sub/models"
	"github.com/gin-gonic/gin"
)

var transactions = []models.Transaction{
	{
		ID:            1,
		Amount:        100,
		AccountName:   "John Doe",
		AccountNumber: "1234567890",
		Direction:     "in",
		Status:        "pending",
		Reference:     "1234567890",
	},
	{
		ID:            2,
		Amount:        200,
		AccountName:   "Jane Doe",
		AccountNumber: "1234567764",
		Direction:     "out",
		Status:        "pending",
		Reference:     "1234567764",
	},
	{
		ID:            3,
		Amount:        300,
		AccountName:   "Jason Deen",
		AccountNumber: "1234567123",
		Direction:     "in",
		Status:        "pending",
		Reference:     "1234567123",
	},
}

func (s *Server) HealthCheck(c *Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) GetTransactions(c *Context) {
	c.JSON(http.StatusOK, transactions)
}

func (s *Server) GetTransaction(c *Context) {
	id := c.Param("id")

	transactionID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	for _, transaction := range transactions {
		if transaction.ID == transactionID {
			c.JSON(http.StatusOK, transaction)
			return
		}
	}

	// Transaction not found
	c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
}
