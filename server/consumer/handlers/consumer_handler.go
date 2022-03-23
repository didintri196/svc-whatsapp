package handlers

import (
	"encoding/json"
	"fmt"
	"svc-whatsapp/domain/payload"
	"svc-whatsapp/usecase"
)

type WhatsappConsumerHandler struct {
	Handler
}

func NewWhatsappConsumerHandler(handler Handler) WhatsappConsumerHandler {
	return WhatsappConsumerHandler{handler}
}

func (handler WhatsappConsumerHandler) ConsumeSendMessage(message []byte) {
	// input & validate json - start //
	input := new(payload.SendMessagePayload)
	if err := json.Unmarshal(message, &input); err != nil {
		panic(err)
	}
	if err := handler.Contract.Validator.Struct(input); err != nil {
		return
	}
	// input & validate json - end //

	// database proccesing - start //
	fmt.Println("TO: ", input.To)
	fmt.Println("MESSAGE: ", input.Message)
	uc := usecase.NewWhatsappWorkerUsecase(handler.Contract)
	err := uc.SendMessage(input)
	if err != nil {
		panic(err)
	}
	// database proccesing - end //

}

func (handler WhatsappConsumerHandler) RegisterWorkerKeeper() {
	go usecase.NewWhatsappKeeperUsecase(handler.Contract).Respawn()
}
