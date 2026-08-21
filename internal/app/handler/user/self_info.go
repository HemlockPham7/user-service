package user

import (
	"net/http"

	"github.com/HemlockPham7/common-libs/pkg/requestutils"
	"github.com/HemlockPham7/common-libs/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

// GetSelfInfo get your current information
// @Summary get your current information
// @Description get your current information
// @Tags User
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} object{data=model.User} "Success"
// @Router /v1/self/info [get]
func (h *userHandler) GetSelfInfo(c *gin.Context) {
	span := newrelic.FromContext(c).StartSegment("GetSelfInfo_UserHandler")
	defer span.End()

	uid, err := requestutils.GetUserIDFromRequest(c)
	if err != nil {
		return
	}

	currentUser, err := h.service.GetSelfInfo(c, uid)
	if err != nil {
		log.Err(err).Msg("Failed to get self info")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, currentUser)
}
