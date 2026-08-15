package user

import (
	"errors"
	"net/http"

	"github.com/HemlockPham7/common-libs/pkg/dbutils"
	"github.com/HemlockPham7/common-libs/pkg/requestutils"
	"github.com/HemlockPham7/common-libs/pkg/response"
	"github.com/HemlockPham7/user-service/internal/app/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type registerInput struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,gte=8"`
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
}

type registerResponse struct {
	Data    *model.User `json:"data"`
	Message string      `json:"message"`
}

// Register handles user registration requests.
// It validates the JSON input, delegates to the service layer for user creation,
// and returns the created user or an appropriate error response.
//
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags User
// @Accept json
// @Produce json
// @Param user body registerInput true "User registration details"
// @Success 201 {object} registerResponse
// @Router /v1/users/register [post]
func (h *userHandler) Register(c *gin.Context) {
	input, err := requestutils.BindInputFromRequest[registerInput](c)
	if err != nil {
		return
	}

	createdUser, err := h.service.CreateUser(c, input.Username, input.Password, input.DisplayName, input.Email)
	switch {
	case errors.Is(err, dbutils.ErrDuplicationUsername):
		c.AbortWithStatusJSON(http.StatusConflict, response.Message{
			Message: "Username already taken",
		})
		return
	case errors.Is(err, dbutils.ErrDuplicationEmail):
		c.AbortWithStatusJSON(http.StatusConflict, response.Message{
			Message: "Email already taken",
		})
		return
	case err == nil:
	default:
		log.Err(err).Msg("Failed to register user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusCreated, &registerResponse{
		Data:    createdUser,
		Message: "Register an user successfully!",
	})
}
