package handlers

import "github.com/gin-gonic/gin"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.AbortWithStatusJSON(status, Response{Status: status, Message: message, Data: data})
}

func NewErrorResponse(ctx *gin.Context, status int, message string) {
	ctx.AbortWithStatusJSON(status, Response{Status: status, Message: message})
}
