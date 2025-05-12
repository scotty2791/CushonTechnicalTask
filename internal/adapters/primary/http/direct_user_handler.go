package http

import (
	"net/http"

	"cushon/internal/core/domain"
	"cushon/internal/core/ports/input"

	"github.com/gin-gonic/gin"
)

// DirectUserHandler handles HTTP requests for direct user operations
type DirectUserHandler struct {
	directUserService input.DirectUserService
}

// NewDirectUserHandler creates a new direct user handler
func NewDirectUserHandler(directUserService input.DirectUserService) *DirectUserHandler {
	return &DirectUserHandler{
		directUserService: directUserService,
	}
}

// RegisterRoutes registers the direct user routes
func (h *DirectUserHandler) RegisterRoutes(router *gin.Engine) {
	directUsers := router.Group("/direct-users")
	{
		directUsers.POST("", h.CreateDirectUser)
		directUsers.GET("/:id", h.GetDirectUser)
		directUsers.PUT("/:id", h.UpdateDirectUser)
		directUsers.DELETE("/:id", h.DeleteDirectUser)
	}
}

// CreateDirectUser handles direct user creation
func (h *DirectUserHandler) CreateDirectUser(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	directUser, err := h.directUserService.CreateDirectUser(request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, directUser)
}

// GetDirectUser handles direct user retrieval
func (h *DirectUserHandler) GetDirectUser(c *gin.Context) {
	id := c.Param("id")
	directUser, err := h.directUserService.GetDirectUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "direct user not found"})
		return
	}

	c.JSON(http.StatusOK, directUser)
}

// UpdateDirectUser handles direct user updates
func (h *DirectUserHandler) UpdateDirectUser(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	directUser := &domain.DirectUser{
		ID:   id,
		Name: request.Name,
	}

	if err := h.directUserService.UpdateDirectUser(directUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, directUser)
}

// DeleteDirectUser handles direct user deletion
func (h *DirectUserHandler) DeleteDirectUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.directUserService.DeleteDirectUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
} 