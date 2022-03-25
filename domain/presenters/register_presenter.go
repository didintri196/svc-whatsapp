package presenters

import (
	"svc-whatsapp/domain/models"
)

type RegisterPresenter struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Jwt   string `json:"jwt"`
}

func NewRegisterPresenter() RegisterPresenter {
	return RegisterPresenter{}
}

func (presenter RegisterPresenter) Build(model *models.MUser, jwt string) RegisterPresenter {
	return RegisterPresenter{
		ID:    model.ID,
		Name:  model.Name,
		Email: model.Email,
		Jwt:   jwt,
	}
}
