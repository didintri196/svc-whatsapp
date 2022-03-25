package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"svc-whatsapp/domain/constants"
	"svc-whatsapp/domain/constants/messages"
	"svc-whatsapp/libraries"
	handler "svc-whatsapp/server/http/handlers/helper"
	"svc-whatsapp/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTMiddleware struct {
	*usecase.Contract
}

func NewJWTMiddleware(ucContract *usecase.Contract) JWTMiddleware {
	return JWTMiddleware{ucContract}
}

func (middleware JWTMiddleware) LoginOnly(ctx *gin.Context) {

	// validate
	if err := middleware.validate(ctx); err != nil {
		handler.NewHandler(middleware.Contract).SendResponseWithoutMeta(ctx, err.Error(), nil, http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (middleware JWTMiddleware) validate(ctx *gin.Context) (err error) {

	// check header
	header := ctx.GetHeader("Authorization")
	if !strings.Contains(header, "Bearer") {
		return errors.New(messages.TokenIsNotProvidedMessage)
	}

	// check token is valid
	token := strings.Replace(header, "Bearer ", "", -1)
	claims, IsValid := libraries.NewJWTLibrary(middleware.SecretKey).ValidateToken(token)
	if !IsValid {
		return errors.New(messages.TokenIsNotValidMessage)
	}

	// check live time
	if expInt, ok := claims[constants.JWTPayloadTokenLiveTime].(float64); ok {
		now := time.Now().Unix()
		if now > int64(expInt) {
			return errors.New(messages.TokenIsExpiredMessage)
		}
	} else {
		return errors.New(messages.InterfaceConversionErrorMessage)
	}

	// insert payload
	if err = middleware.insertPayload(ctx, claims); err != nil {
		return err
	}

	return err
}

func (middleware JWTMiddleware) insertPayload(ctx *gin.Context, claims jwt.MapClaims) (err error) {

	// UUID
	if payloadID, ok := claims[constants.JWTPayloadUUID].(string); ok {
		ctx.Set(constants.JWTPayloadUUID, payloadID)
	} else {
		return errors.New(messages.InterfaceConversionErrorMessage)
	}

	// EMAIL
	if payloadEmail, ok := claims[constants.JWTPayloadEmail].(string); ok {
		ctx.Set(constants.JWTPayloadEmail, payloadEmail)
	} else {
		return errors.New(messages.InterfaceConversionErrorMessage)
	}

	return err
}
