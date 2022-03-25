package usecase

import (
	"fmt"
	"svc-whatsapp/domain/payload"
	"svc-whatsapp/libraries"
)

type (
	IWhatsappWorkerUsecase interface {
		SendMessage(payload *payload.SendMessagePayload) (err error)
	}

	WhatsappWorkerUsecase struct {
		*Contract
	}
)

func NewWhatsappWorkerUsecase(ucContract *Contract) IWhatsappWorkerUsecase {
	return &WhatsappWorkerUsecase{ucContract}
}

func (uc WhatsappWorkerUsecase) SendMessage(payload *payload.SendMessagePayload) (err error) {
	fmt.Println(uc.WhatsappWorker.GetAllIdle())
	uc.WhatsappWorker.Publish("000-0", libraries.SendMessage{
		To:      payload.Message,
		Message: payload.Message,
	})
	return err
}
