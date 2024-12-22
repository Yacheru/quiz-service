package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quiz-service/internal/server/http/handlers"
)

func VariantName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		variantName := ctx.Param("variantName")
		if variantName == "" {
			handlers.NewErrorResponse(ctx, http.StatusBadRequest, "no variant name supplied")
		}

		ctx.Next()
	}
}
