package handler

import (
	"svc-whatsapp/domain/presenters"
	"svc-whatsapp/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Contract *usecase.Contract
}

func NewHandler(ucContract *usecase.Contract) Handler {
	return Handler{Contract: ucContract}
}

func (handler Handler) SendResponseWithoutMeta(ctx *gin.Context, message string, data interface{}, httpStatus int) {
	ctx.JSON(httpStatus, presenters.ResponsePresenter{
		Message: message,
		Data:    data,
	})
}

func (handler Handler) SendResponseWithMeta(ctx *gin.Context, data interface{}, message string, meta interface{}, httpStatus int) {
	ctx.JSON(httpStatus, presenters.ResponsePresenter{
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
