package handlers

import (
	"net/http"
	"svc-whatsapp/domain/constants/messages"
	"svc-whatsapp/domain/requests"
	handler "svc-whatsapp/server/http/handlers/helper"
	"svc-whatsapp/usecase"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	handler.Handler
}

func NewDeviceHandler(handler handler.Handler) DeviceHandler {
	return DeviceHandler{handler}
}

func (handler DeviceHandler) Filter(ctx *gin.Context) {
	// input & validate json - start //
	input := new(requests.FilterRequest)
	if err := ctx.Bind(&input); err != nil {
		handler.SendResponseWithoutMeta(ctx, messages.FailedMessage, nil, http.StatusBadRequest)
		return
	}
	if err := handler.Contract.Validator.Struct(input); err != nil {
		handler.SendResponseWithoutMeta(ctx, err.Error(), nil, http.StatusBadRequest)
		return
	}
	// input & validate json - end //

	// database proccesing - start //
	uc := usecase.NewMDevicesUsecase(handler.Contract)
	presenter, meta, err := uc.FilterMDevices(ctx, input)
	if err != nil {
		handler.SendResponseWithoutMeta(ctx, messages.FailedMessage, nil, http.StatusBadRequest)
		return
	}

	// database proccesing - end //
	handler.SendResponseWithMeta(ctx, presenter.FilterDevicesPresenter, messages.SuccessMessage, meta, http.StatusOK)
}
