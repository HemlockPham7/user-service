package user

import (
	"github.com/HemlockPham7/user-service/internal/app/service/user"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetSelfInfo(c *gin.Context)
	UpdateUserByID(c *gin.Context)
}

type userHandler struct {
	service user.Service
}

func NewHandler(service user.Service) Handler {
	return &userHandler{service: service}
}
