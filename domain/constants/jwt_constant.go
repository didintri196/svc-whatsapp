package constants

import "time"

const (
	JWTTokenLiveTIme        = 24 * time.Hour
	JWTTokenOtpLiveTIme     = 15 * time.Minute
	JWTRefreshTokenLiveTime = 24 * time.Hour

	JWTResponseToken        = "token"
	JWTResponseRefreshToken = "refresh_token"

	JWTPayloadUUID          = "uuid"
	JWTPayloadEmail         = "email"
	JWTPayloadRoleName      = "role_name"
	JWTPayloadOtpName       = "otp"
	JWTPayloadTokenLiveTime = "exp"
)
