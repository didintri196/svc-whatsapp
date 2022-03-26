package usecase

import (
	"fmt"
	"svc-whatsapp/domain/models"
	"svc-whatsapp/domain/payload"
	"svc-whatsapp/libraries"
	"svc-whatsapp/repositories"
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

	repo := repositories.NewMDevicesRepository(uc.Postgres)
	modelDevice := models.NewMDevices()
	err = repo.ReadByApiKey(payload.ApiKey, modelDevice)

	if modelDevice.WorkerID != "" {
		fmt.Println("FOUND WORKER ID " + modelDevice.WorkerID)
		uc.WhatsappWorker.Publish(modelDevice.WorkerID, libraries.SendMessage{
			To:      payload.To,
			Message: payload.Message,
		})
	} else {
		//todo retry queue nsq until found worker id
	}
	return err
}
