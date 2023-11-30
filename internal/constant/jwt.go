package constant

import "github.com/golang-jwt/jwt/v4"

var (
	JWTSigningMethod = jwt.SigningMethodHS256
)

type TokenType string

const (
	AccessTokenType  TokenType = "AccessToken"
	RefreshTokenType TokenType = "RefreshToken"
	StepUpTokenType  TokenType = "StepUpToken"
)
