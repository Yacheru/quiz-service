package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"quiz-service/init/logger"
	"quiz-service/internal/entities"
	"quiz-service/internal/service"
	"quiz-service/pkg/constants"
)

type Handler struct {
	service *service.Service

	logger logger.Logging
}

func NewHandler(service *service.Service, logger logger.Logging) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Register(ctx *gin.Context) {
	h.logger.InfoF("register handler received by: %s", ctx.Request.UserAgent())

	registerEntity := new(entities.Register)
	if err := ctx.ShouldBindBodyWithJSON(registerEntity); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.RegisterService.Register(ctx.Request.Context(), registerEntity)
	if err != nil {
		if errors.Is(err, constants.ErrorUserAlreadyExists) {
			NewErrorResponse(ctx, http.StatusConflict, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusCreated, "User successfully created", user)
	return
}

func (h *Handler) Login(ctx *gin.Context) {
	h.logger.InfoF("login handler received by: %s", ctx.Request.UserAgent())

	loginEntity := new(entities.Login)
	if err := ctx.ShouldBindBodyWithJSON(loginEntity); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.RegisterService.Login(ctx.Request.Context(), loginEntity)
	if err != nil {
		if errors.Is(err, constants.ErrorUserNotFound) {
			NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusCreated, "User successfully login", user)
	return
}
