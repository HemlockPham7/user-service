package user

import (
	"errors"
	"net/http"

	"github.com/HemlockPham7/common-libs/pkg/requestutils"
	"github.com/HemlockPham7/common-libs/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type updateUserRequest struct {
	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email,omitempty"`
}

// UpdateUserByID updates user by ID
// @Summary update user by ID
// @Description update user by ID
// @Tags User
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Param user body updateUserRequest true "User update details"
// @Success 200 {object} response.Message
// @Router /v1/users/update [put]
func (h *userHandler) UpdateUserByID(c *gin.Context) {
	request, uid, err := requestutils.BindInputFromRequestWithAuth[updateUserRequest](c)
	if err != nil {
		return
	}

	err = h.service.UpdateUserByID(c, uid, request.DisplayName, request.Email)
	switch {
	case errors.Is(err, nil):
		break
	default:
		log.Err(err).Str("operation", "UpdateUser").Msg("service return error when update user")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, response.Message{Message: "User updated successfully"})
}
