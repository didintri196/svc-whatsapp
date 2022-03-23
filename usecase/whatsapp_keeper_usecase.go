package usecase

import "svc-whatsapp/libraries"

type (
	IWhatsappKeeperUsecase interface {
		Respawn()``
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
		idWorker := uc.WhatsappWorker.GetOneIdle()
		if idWorker != "" {
			//TODO get in db jid not set worker
			// And set data to publish worker
			jdid := "6285155075517.0:11@s.whatsapp.net"
			// publish to worker
			uc.WhatsappWorker.Publish(idWorker, libraries.ConnectMessage{JDID: jdid})
		}
	}

}
