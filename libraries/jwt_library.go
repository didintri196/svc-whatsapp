package libraries

import (
	"errors"
	"svc-whatsapp/domain/constants"
	"svc-whatsapp/domain/constants/messages"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTLibrary struct {
	secretKey []byte
}

func NewJWTLibrary(secretKey []byte) JWTLibrary {
	return JWTLibrary{
		secretKey: secretKey,
	}
}

func (lib JWTLibrary) GenerateToken(id, email string) (jwt string, err error) {

	//token
	token, err := lib.generateTokenWithTTL(id, email, constants.JWTTokenLiveTIme)
	if err != nil {
		return jwt, err
	}

	return token, err
}

func (lib JWTLibrary) ValidateToken(encodedToken string) (jwt.MapClaims, bool) {

	// parse token
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return encodedToken, errors.New(messages.TokenIsNotValidMessage)
		}
		return lib.secretKey, nil
	})
	if err != nil {
		return nil, false
	}

	// get payload
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	}

	return nil, false
}

func (lib JWTLibrary) generateTokenWithTTL(id, email string, duration time.Duration) (signedToken string, err error) {

	// payload
	payload := jwt.MapClaims{}
	payload[constants.JWTPayloadUUID] = id
	payload[constants.JWTPayloadEmail] = email
	payload[constants.JWTPayloadTokenLiveTime] = time.Now().Add(duration).Unix()

	// token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload) // Encryption Algorithm

	// signature
	signedToken, err = token.SignedString(lib.secretKey)
	if err != nil {
		return signedToken, err
	}

	return signedToken, err
}

func (lib JWTLibrary) GenerateTokenOtp(id, email, otp string) (token string, err error) {
	//token
	token, err = lib.generateTokenOtpWithTTL(id, email, otp, constants.JWTTokenOtpLiveTIme)
	if err != nil {
		return token, err
	}

	return token, err
}

func (lib JWTLibrary) generateTokenOtpWithTTL(id, email, otp string, duration time.Duration) (signedToken string, err error) {

	// payload
	payload := jwt.MapClaims{}
	payload[constants.JWTPayloadUUID] = id
	payload[constants.JWTPayloadEmail] = email
	payload[constants.JWTPayloadOtpName] = otp
	payload[constants.JWTPayloadTokenLiveTime] = time.Now().Add(duration).Unix()

	// token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload) // Encryption Algorithm

	// signature
	signedToken, err = token.SignedString(lib.secretKey)
	if err != nil {
		return signedToken, err
	}

	return signedToken, err
}
