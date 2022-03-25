package presenters

import (
	"svc-whatsapp/domain/models"
	"time"
)

type ArrayFilterDevicesPresenter struct {
	FilterDevicesPresenter []FilterDevicesPresenter
}

type FilterDevicesPresenter struct {
	ID        string    `json:"id"`
	Jid       string    `json:"jid"`
	Server    string    `json:"server"`
	Phone     string    `json:"phone"`
	WorkerID  string    `json:"worker_id"`
	ApiKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewArrayFilterDevicesPresenter() ArrayFilterDevicesPresenter {
	return ArrayFilterDevicesPresenter{}
}

func (presenter ArrayFilterDevicesPresenter) Build(model []models.Devices) (list ArrayFilterDevicesPresenter) {
	for _, row := range model {
		list.FilterDevicesPresenter = append(list.FilterDevicesPresenter, FilterDevicesPresenter{
			ID:        row.ID,
			Jid:       row.Jid,
			Server:    row.Server,
			Phone:     row.Phone,
			WorkerID:  row.WorkerID,
			ApiKey:    row.ApiKey,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		})
	}
	return list
}
