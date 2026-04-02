package handler

import "github.com/gin-gonic/gin"

// envelope is the standard response wrapper: { "data": ..., "errors": [...] }
type envelope struct {
	Data   any   `json:"data"`
	Errors []any `json:"errors,omitempty"`
}

// respondData sends a 200 response with data in the standard envelope.
func respondData(c *gin.Context, data any) {
	c.JSON(200, envelope{Data: data})
}

// respondError sends an error response in the standard envelope.
func respondError(c *gin.Context, status int, msg string) {
	c.JSON(status, envelope{
		Errors: []any{gin.H{"message": msg}},
	})
}
