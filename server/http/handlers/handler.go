package handlers

import (
	"encoding/json"
	"net/http"
	"svc-whatsapp/domain/constants"
	"svc-whatsapp/domain/constants/messages"
	"svc-whatsapp/domain/payload"
	"svc-whatsapp/domain/requests"
	handler "svc-whatsapp/server/http/handlers/helper"

	"github.com/gin-gonic/gin"
)

type WhatsappApiHandler struct {
	handler.Handler
}

func NewWhatsappApiHandler(handler handler.Handler) WhatsappApiHandler {
	return WhatsappApiHandler{handler}
}

func (handler WhatsappApiHandler) TestProduce(ctx *gin.Context) {
	// input & validate json - start //
	input := new(requests.ProduceReq)
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
	exampleBytes, err := json.Marshal(input.Body)
	if err != nil {
		print(err)
		return
	}
	handler.Contract.NsqProducer.Publish(input.Topic, exampleBytes)
	// database proccesing - end //

	handler.SendResponseWithoutMeta(ctx, messages.SuccessMessage, input.Body, http.StatusOK)
}

func (handler WhatsappApiHandler) SendMessage(ctx *gin.Context) {
	apiKey, found := ctx.GetQuery("api")
	if !found {
		handler.SendResponseWithoutMeta(ctx, messages.FailedMessage, nil, http.StatusBadRequest)
	}
	// input & validate json - start //
	input := new(requests.SendMessageRequest)
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
	exampleBytes, err := json.Marshal(payload.SendMessagePayload{
		ApiKey:  apiKey,
		To:      input.To,
		Message: input.Message,
	})
	if err != nil {
		handler.SendResponseWithoutMeta(ctx, err.Error(), nil, http.StatusBadRequest)
		return
	}
	handler.Contract.NsqProducer.Publish(constants.TopicSendMessage, exampleBytes)
	// database proccesing - end //

	handler.SendResponseWithoutMeta(ctx, messages.QueueMessage, input, http.StatusOK)
}
