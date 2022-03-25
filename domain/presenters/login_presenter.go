package presenters

import (
	"svc-whatsapp/domain/models"
)

type LoginPresenter struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Jwt   string `json:"jwt,omitempty"`
}

func NewLoginPresenter() LoginPresenter {
	return LoginPresenter{}
}

func (presenter LoginPresenter) Build(model *models.MUser, jwt string) LoginPresenter {
	return LoginPresenter{
		ID:    model.ID,
		Name:  model.Name,
		Email: model.Email,
		Jwt:   jwt,
	}
}
