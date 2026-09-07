package user

import (
	"github.com/HemlockPham7/user-service/internal/app/service/user"
	"github.com/gin-gonic/gin"
)

// Handler defines the HTTP handlers for user management operations.
type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetSelfInfo(c *gin.Context)
	UpdateUserByID(c *gin.Context)
}

type userHandler struct {
	service user.Service
}

// NewHandler creates a user HTTP handler with the provided user service.
//
// Parameters:
//   - service: the user service used to execute user management operations.
//
// Returns:
//   - A user HTTP handler configured with the provided service.
func NewHandler(service user.Service) Handler {
	return &userHandler{service: service}
}
