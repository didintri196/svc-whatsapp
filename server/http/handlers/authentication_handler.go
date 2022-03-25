package handlers

import (
	"net/http"
	"svc-whatsapp/domain/constants"
	"svc-whatsapp/domain/constants/messages"
	"svc-whatsapp/domain/presenters"
	"svc-whatsapp/domain/requests"
	handler "svc-whatsapp/server/http/handlers/helper"
	"svc-whatsapp/usecase"

	"github.com/gin-gonic/gin"
)

type AuthenticationHandler struct {
	handler.Handler
}

func NewAuthenticationHandler(handler handler.Handler) AuthenticationHandler {
	return AuthenticationHandler{handler}
}

func (handler AuthenticationHandler) Register(ctx *gin.Context) {
	// input & validate json - start //
	input := new(requests.RegisterRequest)
	if err := ctx.Bind(&input); err != nil {
		handler.SendResponseWithoutMeta(ctx, messages.FailedMessage, nil, http.StatusBadRequest)
		return
	}
	if err := handler.Contract.Validator.Struct(input); err != nil {
		handler.SendResponseWithoutMeta(ctx, err.Error(), nil, http.StatusBadRequest)
		return
	}
	// input & validate json - end //

	// database proccesing - start //
	handler.Contract.PostgresTX = handler.Contract.Postgres.Begin()
	uc := usecase.NewAuthenticationUsecase(handler.Contract)
	res, err := uc.Register(input)
	if err != nil {
		handler.Contract.PostgresTX.Rollback()
		handler.SendResponseWithoutMeta(ctx, err.Error(), nil, http.StatusUnprocessableEntity)
		return
	}
	handler.Contract.PostgresTX.Commit()
	// database proccesing - end //

	handler.SendResponseWithoutMeta(ctx, messages.SuccessMessage, res, http.StatusOK)
}

func (handler AuthenticationHandler) Login(ctx *gin.Context) {
	// input & validate json - start //
	input := new(requests.LoginRequest)
	if err := ctx.Bind(&input); err != nil {
		handler.SendResponseWithoutMeta(ctx, messages.FailedMessage, nil, http.StatusBadRequest)
		return
	}
	if err := handler.Contract.Validator.Struct(input); err != nil {
		handler.SendResponseWithoutMeta(ctx, err.Error(), nil, http.StatusBadRequest)
		return
	}
	// input & validate json - end //

	// database proccesing - start //
	uc := usecase.NewAuthenticationUsecase(handler.Contract)
	res, err := uc.Login(input)
	if err != nil {
		handler.SendResponseWithoutMeta(ctx, err.Error(), nil, http.StatusUnauthorized)
		return
	}
	// database proccesing - end //

	handler.SendResponseWithoutMeta(ctx, messages.SuccessMessage, res, http.StatusOK)
}

func (handler AuthenticationHandler) GetJwt(ctx *gin.Context) {
	id, ok := ctx.Get(constants.JWTPayloadUUID)
	if !ok {
		handler.SendResponseWithoutMeta(ctx, messages.InterfaceConversionErrorMessage, id, http.StatusOK)
		return
	}

	email, ok := ctx.Get(constants.JWTPayloadEmail)
	if !ok {
		handler.SendResponseWithoutMeta(ctx, messages.InterfaceConversionErrorMessage, email, http.StatusOK)
		return
	}

	if ok {
		handler.SendResponseWithoutMeta(ctx, messages.SuccessMessage, presenters.LoginPresenter{
			ID:    id.(string),
			Email: email.(string),
		}, http.StatusOK)
	}
}
