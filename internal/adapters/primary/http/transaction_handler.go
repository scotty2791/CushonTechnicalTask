package http

import (
	"net/http"

	"cushon/internal/core/domain"
	"cushon/internal/core/ports/input"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// TransactionHandler handles HTTP requests for transaction operations
type TransactionHandler struct {
	transactionService input.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService input.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// RegisterRoutes registers the transaction routes
func (h *TransactionHandler) RegisterRoutes(router *gin.Engine) {
	transactions := router.Group("/transactions")
	{
		transactions.POST("", h.CreateTransaction)
		transactions.GET("/:id", h.GetTransaction)
		transactions.GET("/user/:userID", h.GetUserTransactions)
		transactions.PUT("/:id", h.UpdateTransaction)
		transactions.DELETE("/:id", h.DeleteTransaction)
	}

	// Add route for getting available fund names
	router.GET("/fund-names", h.GetFundNames)
}

// GetFundNames returns a list of available fund names
func (h *TransactionHandler) GetFundNames(c *gin.Context) {
	fundNames := []string{
		string(domain.CushonEquitiesFund),
	}
	c.JSON(http.StatusOK, fundNames)
}

// CreateTransaction handles transaction creation
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var request struct {
		UserID   string          `json:"user_id" binding:"required"`
		Amount   decimal.Decimal `json:"amount" binding:"required"`
		FundName string          `json:"fund_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.CreateTransaction(
		request.UserID,
		request.Amount,
		domain.FundName(request.FundName),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetTransaction handles transaction retrieval
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id := c.Param("id")
	transaction, err := h.transactionService.GetTransaction(id)
	if err != nil {
		switch err.Error() {
		case "transaction ID is required":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "transaction not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetUserTransactions handles user transactions retrieval
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	userID := c.Param("userID")
	transactions, err := h.transactionService.GetUserTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// UpdateTransaction handles transaction updates
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Amount   decimal.Decimal `json:"amount" binding:"required"`
		FundName string          `json:"fund_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.GetTransaction(id)
	if err != nil {
		switch err.Error() {
		case "transaction ID is required":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "transaction not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	transaction.Amount = request.Amount
	transaction.FundName = domain.FundName(request.FundName)

	if err := h.transactionService.UpdateTransaction(transaction); err != nil {
		switch err.Error() {
		case "transaction not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "invalid fund name":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction handles transaction deletion
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	if err := h.transactionService.DeleteTransaction(id); err != nil {
		switch err.Error() {
		case "transaction ID is required":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "transaction not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
} 