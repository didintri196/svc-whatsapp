package middlewares

import (
	"errors"
	"net/http"
	"svc-whatsapp/domain/constants/messages"
	handler "svc-whatsapp/server/http/handlers/helper"
	"svc-whatsapp/usecase"

	"github.com/gin-gonic/gin"
)

type QueryMiddleware struct {
	*usecase.Contract
}

func NewQueryiddleware(ucContract *usecase.Contract) QueryMiddleware {
	return QueryMiddleware{ucContract}
}

func (middleware QueryMiddleware) TokenisFilled(ctx *gin.Context) {

	// validate
	if err := middleware.validate(ctx); err != nil {
		handler.NewHandler(middleware.Contract).SendResponseWithoutMeta(ctx, err.Error(), nil, http.StatusUnauthorized)
		ctx.Abort()
		return
	}
	ctx.Next()
}

func (middleware QueryMiddleware) validate(ctx *gin.Context) (err error) {

	// check header
	ctx.Request.Header.Del("Origin")
	token, _ := ctx.GetQuery("hex")
	if token == "" {
		return errors.New(messages.TokenIsNotProvidedMessage)
	}

	ctx.Set("ID", token)
	return err
}
