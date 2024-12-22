package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"quiz-service/internal/entities"
	"quiz-service/pkg/constants"
)

func (h *Handler) VariantAdd(ctx *gin.Context) {
	h.logger.InfoF("VariantAdd handler received by: %s", ctx.Request.UserAgent())

	questionEntity := new(entities.Variant)
	if err := ctx.ShouldBindBodyWithJSON(questionEntity); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.VariantService.VariantAdd(ctx.Request.Context(), questionEntity.Name); err != nil {
		if errors.Is(err, constants.ErrorVariantAlreadyExists) {
			NewErrorResponse(ctx, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, constants.ErrorVariantTooLong) {
			NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusCreated, "Variant successfully created", nil)
	return
}

func (h *Handler) VariantRemove(ctx *gin.Context) {
	h.logger.InfoF("VariantRemove handler received by: %s", ctx.Request.UserAgent())

	variant := ctx.MustGet("variant").(*entities.Variant)

	if err := h.service.VariantService.VariantRemove(ctx.Request.Context(), variant.Name); err != nil {
		if errors.Is(err, constants.ErrorVariantNotFound) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "Variant successfully removed", nil)
	return
}

func (h *Handler) VariantList(ctx *gin.Context) {
	h.logger.InfoF("VariantList handler received by: %s", ctx.Request.UserAgent())

	variants, err := h.service.VariantService.VariantList(ctx.Request.Context())
	if err != nil {
		if errors.Is(err, constants.ErrorNoVariantsYet) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "all variants", variants)
	return
}

func (h *Handler) VariantCheck(ctx *gin.Context) {
	h.logger.InfoF("VariantCheck handler received by: %s", ctx.Request.UserAgent())

	variantName := ctx.Param("variantName")
	variant, err := h.service.VariantService.VariantGet(ctx.Request.Context(), variantName)
	if err != nil {
		if errors.Is(err, constants.ErrorVariantNotFound) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Set("variant", variant)
	ctx.Next()
}

func (h *Handler) VariantGet(ctx *gin.Context) {
	h.logger.InfoF("VariantGet handler received by: %s", ctx.Request.UserAgent())

	NewSuccessResponse(ctx, http.StatusOK, "variant", ctx.MustGet("variant").(*entities.Variant))
	return
}

func (h *Handler) VariantStart(ctx *gin.Context) {
	h.logger.InfoF("VariantStart handler received by: %s", ctx.Request.UserAgent())

	variant := ctx.MustGet("variant").(*entities.Variant)
	user := ctx.MustGet("user").(*entities.User)

	if err := h.service.VariantService.VariantStart(ctx.Request.Context(), variant.Id, user.ID); err != nil {
		if errors.Is(err, constants.ErrorVariantNotFound) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, constants.ErrorVariantCompleted) {
			NewErrorResponse(ctx, http.StatusConflict, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "variant successfully started", variant)
	return
}

func (h *Handler) VariantResults(ctx *gin.Context) {
	h.logger.InfoF("VariantResults handler received by: %s", ctx.Request.UserAgent())

	variant := ctx.MustGet("variant").(*entities.Variant)
	user := ctx.MustGet("user").(*entities.User)

	testing, err := h.service.VariantService.VariantResults(ctx.Request.Context(), variant.Id, user.ID)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "results.html", gin.H{
		"correctAnswers": testing.CorrectAnswers,
	})
	return
}
