package usecase

import (
	"errors"
	"svc-whatsapp/domain/constants/messages"
	"svc-whatsapp/domain/models"
	"svc-whatsapp/domain/presenters"
	"svc-whatsapp/domain/requests"
	"svc-whatsapp/libraries"
	"svc-whatsapp/repositories"
	"svc-whatsapp/utils"
	"time"

	"github.com/google/uuid"
)

type (
	IAuthenticationUsecase interface {
		Register(input *requests.RegisterRequest) (presenter presenters.RegisterPresenter, err error)
		Login(input *requests.LoginRequest) (presenter presenters.LoginPresenter, err error)
	}

	AuthenticationUsecase struct {
		*Contract
	}
)

func NewAuthenticationUsecase(ucContract *Contract) IAuthenticationUsecase {
	return &AuthenticationUsecase{ucContract}
}

func (uc AuthenticationUsecase) Register(input *requests.RegisterRequest) (presenter presenters.RegisterPresenter, err error) {

	// init
	now := time.Now()
	password, err := utils.NewHashHelper().HashAndSalt(input.Password)
	if err != nil {
		NewErrorLog("AuthenticationUsecase.Register", "utils.HashAndSalt", err.Error())
		return presenter, err
	}
	id := uuid.New().String()
	model := &models.MUser{
		ID:        id,
		Name:      input.Name,
		Email:     input.Email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save not verified user
	repo := repositories.NewMUserRepository(uc.Postgres)
	if err = repo.Create(uc.PostgresTX, model); err != nil {
		NewErrorLog("AuthenticationUseCase.Register", "repo.Create", err.Error())
		return presenter, err
	}

	// generate jwt
	jwtLibrary := libraries.NewJWTLibrary(uc.SecretKey)
	jwt, err := jwtLibrary.GenerateToken(model.ID, model.Email)
	if err != nil {
		NewErrorLog("AuthenticationUseCase.Register", "jwtLibrary.GenerateToken", err.Error())
		return presenter, err
	}

	// build presenter
	presenter = presenters.NewRegisterPresenter().Build(model, jwt)

	return presenter, err
}

func (uc AuthenticationUsecase) Login(input *requests.LoginRequest) (presenter presenters.LoginPresenter, err error) {

	// get user model
	model := models.NewMUser()
	repo := repositories.NewMUserRepository(uc.Postgres)
	if err = repo.ReadByEmail(input.Email, model); err != nil {
		NewErrorLog("Authentication.Login", "repo.ReadByEmail", err.Error())
		return presenter, err
	}

	// check password is valid
	if isValid := utils.NewHashHelper().CheckHashString(input.Password, model.Password); !isValid {
		return presenter, errors.New(messages.CredentialIsNotMatchMessage)
	}

	// send jwt
	jwtLibrary := libraries.NewJWTLibrary(uc.SecretKey)
	jwt, err := jwtLibrary.GenerateToken(model.ID, model.Email)
	if err != nil {
		NewErrorLog("Authentication.Login", "jwtLibrary.GenerateToken", err.Error())
		return presenter, err
	}

	// build presenter
	presenter = presenters.NewLoginPresenter().Build(model, jwt)

	return presenter, err
}
