package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quiz-service/internal/server/http/handlers"
	"strconv"
)

func QuestionId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, err := strconv.Atoi(ctx.Param("questionId")); err != nil {
			handlers.NewErrorResponse(ctx, http.StatusBadRequest, "Question id must be a positive integer")
			return
		}

		ctx.Next()
	}
}
