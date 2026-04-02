package handler

import "github.com/gin-gonic/gin"

// Health returns a simple health check response.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
