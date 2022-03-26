package presenters

import (
	"svc-whatsapp/domain/models"
	"time"
)

type DevicesPresenter struct {
	ID        string    `json:"id"`
	Jid       string    `json:"jid"`
	Server    string    `json:"server"`
	Phone     string    `json:"phone"`
	WorkerID  string    `json:"worker_id"`
	ApiKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewDevicesPresenter() DevicesPresenter {
	return DevicesPresenter{}
}

func (presenter DevicesPresenter) Build(model *models.Devices) DevicesPresenter {
	return DevicesPresenter{
		ID:        model.ID,
		Jid:       model.Jid,
		Server:    model.Server,
		Phone:     model.Phone,
		WorkerID:  model.WorkerID,
		ApiKey:    model.ApiKey,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
