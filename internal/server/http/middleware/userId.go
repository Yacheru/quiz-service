package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"quiz-service/internal/server/http/handlers"
)

func UserId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("userId")
		if id == "" {
			handlers.NewErrorResponse(ctx, http.StatusBadRequest, "No id supplied")
			return
		}

		if _, err := uuid.Parse(id); err != nil {
			handlers.NewErrorResponse(ctx, http.StatusBadRequest, "Invalid id supplied")
			return
		}

		ctx.Next()
	}
}
