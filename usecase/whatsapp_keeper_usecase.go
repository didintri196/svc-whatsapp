package usecase

import (
	"log"
	"svc-whatsapp/domain/models"
	"svc-whatsapp/libraries"
	"svc-whatsapp/repositories"
	"time"
)

type (
	IWhatsappKeeperUsecase interface {
		Respawn()
	}

	WhatsappKeeperUsecase struct {
		*Contract
	}
)

func NewWhatsappKeeperUsecase(ucContract *Contract) IWhatsappKeeperUsecase {
	return &WhatsappKeeperUsecase{ucContract}
}

func (uc WhatsappKeeperUsecase) Respawn() {
	for {
		allWorker := uc.WhatsappWorker.GetAllIdle()
		log.Println("[WORKER] Get all worker idle ", allWorker)
		idWorker := uc.WhatsappWorker.GetOneIdle()
		log.Println("[WORKER] Get worker id " + idWorker)
		if idWorker != "" {
			//TODO get in db jid not set worker
			model := models.NewMDevices()
			repo := repositories.NewMDevicesRepository(uc.Postgres)
			if err := repo.ReadNotConnectWorker(model); err != nil {
				NewErrorLog("WhatsappKeeperUsecase.Respawn", "repo.ReadNotConnectWorker", err.Error())
			}
			// And set data to publish worker
			//jdid := "6285155075517.0:11@s.whatsapp.net"
			// publish to worker
			if model.ID != "" {
				uc.WhatsappWorker.Publish(idWorker, libraries.ConnectMessage{JDID: model.Jid})
			}
		}
		time.Sleep(5 * time.Second)
	}

}
