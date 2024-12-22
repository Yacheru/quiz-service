package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"quiz-service/internal/entities"
	"quiz-service/pkg/constants"
)

func (h *Handler) Authenticated(ctx *gin.Context) {
	h.logger.InfoF("Authenticated handler received by: %s", ctx.Request.UserAgent())

	user, err := h.service.UserService.Authenticated(ctx.Request.Context(), ctx.Param("userId"))
	if err != nil {
		if errors.Is(err, constants.ErrorUserNotFound) || errors.Is(err, constants.ErrorUserNotAuthorized) {
			NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}

func (h *Handler) Quit(ctx *gin.Context) {
	h.logger.InfoF("Quit handler received by: %s", ctx.Request.UserAgent())

	user := ctx.MustGet("user").(*entities.User)

	if err := h.service.UserService.Quit(ctx.Request.Context(), user.UUID); err != nil {
		if errors.Is(err, constants.ErrorUserNotFound) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "Quit successfully", nil)
	return
}
