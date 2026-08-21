package user

import (
	"errors"
	"net/http"

	"github.com/HemlockPham7/common-libs/pkg/requestutils"
	"github.com/HemlockPham7/common-libs/pkg/response"
	"github.com/HemlockPham7/user-service/internal/app/service/user"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

type loginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
}

// Login Authentication endpoint
// @Summary Return a jwt token if the input is correct
// @Description Return a jwt token if the input is correct
// @Tags User
// @Accept application/json
// @Produce application/json
// @Param input body loginInput true "Input required"
// @Success 200 {object} object{token=string,message=string} "Success"
// @Router /v1/users/login [post]
func (h *userHandler) Login(c *gin.Context) {
	span := newrelic.FromContext(c).StartSegment("Login_UserHandler")
	defer span.End()

	// doc body input
	input, err := requestutils.BindInputFromRequest[loginInput](c)
	if err != nil {
		return
	}

	// call service -> service tra ve token
	token, err := h.service.Login(c, input.Username, input.Password)
	switch {
	case errors.Is(err, user.ErrInvalidCredentials):
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Message{Message: "invalid credentials"})
		return
	case err == nil:
	default:
		log.Err(err).Msg("Failed to login user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	// tra ve token
	c.JSON(http.StatusOK, &loginResponse{
		Data:    token,
		Message: "Logged in successfully!",
	})
}
