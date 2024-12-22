package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"quiz-service/internal/entities"
	"quiz-service/pkg/constants"
	"strconv"
)

func (h *Handler) QuestionAdd(ctx *gin.Context) {
	h.logger.InfoF("QuestionAdd handler received by: %s", ctx.Request.UserAgent())

	questionEntity := new(entities.Question)
	if err := ctx.ShouldBindBodyWithJSON(questionEntity); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	variant := ctx.MustGet("variant").(*entities.Variant)

	if err := h.service.QuestionsService.QuestionAdd(ctx.Request.Context(), variant.Id, questionEntity); err != nil {
		if errors.Is(err, constants.ErrorVariantNotFound) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, constants.ErrorQuestionLimitExceeded) {
			NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, constants.ErrorQuestionAlreadyExists) {
			NewErrorResponse(ctx, http.StatusConflict, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "Question added successfully", nil)
	return
}

func (h *Handler) QuestionRemove(ctx *gin.Context) {
	h.logger.InfoF("QuestionRemove handler received by: %s", ctx.Request.UserAgent())

	questionEntity := new(entities.QuestionRemove)
	if err := ctx.ShouldBindBodyWithJSON(questionEntity); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	variant := ctx.MustGet("variant").(*entities.Variant)

	if err := h.service.QuestionsService.QuestionRemove(ctx.Request.Context(), variant.Id, questionEntity); err != nil {
		if errors.Is(err, constants.ErrorVariantNotFound) || errors.Is(err, constants.ErrorQuestionNotFound) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "Question removed successfully", nil)
	return
}

func (h *Handler) QuestionGet(ctx *gin.Context) {
	h.logger.InfoF("QuestionGet handler received by: %s", ctx.Request.UserAgent())

	questionId, _ := strconv.Atoi(ctx.Param("questionId"))
	variant := ctx.MustGet("variant").(*entities.Variant)

	question, err := h.service.QuestionsService.QuestionGet(ctx.Request.Context(), variant.Id, questionId)
	if err != nil {
		if errors.Is(err, constants.ErrorQuestionNotFound) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "question", question)
	return
}

func (h *Handler) QuestionAccept(ctx *gin.Context) {
	h.logger.InfoF("QuestionAccept handler received by: %s", ctx.Request.UserAgent())

	answerEntity := new(entities.Answer)
	if err := ctx.ShouldBindBodyWithJSON(answerEntity); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	variant := ctx.MustGet("variant").(*entities.Variant)
	user := ctx.MustGet("user").(*entities.User)

	if err := h.service.QuestionsService.QuestionAccept(ctx.Request.Context(), variant.Id, user.ID, answerEntity.Answer); err != nil {
		if errors.Is(err, constants.ErrorQuestionNotFound) || errors.Is(err, constants.ErrorTestNotFound) {
			NewErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	NewSuccessResponse(ctx, http.StatusOK, "answer accepted", nil)
	return
}
