package shared

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
)

type (
	AccessJWTClaim struct {
		jwt.RegisteredClaims
		UserId    int64              `json:"user_id"`
		IsSeller  bool               `json:"is_seller"`
		TokenType constant.TokenType `json:"token_type"`
	}
	RefreshJWTClaim struct {
		jwt.RegisteredClaims
		TokenType constant.TokenType `json:"token_type"`
	}
	StepUpJWTClaim struct {
		jwt.RegisteredClaims
		TokenType constant.TokenType `json:"token_type"`
	}
	SignAccessTokenPayload struct {
		UserID   int64
		IsSeller bool
	}
)

func (c AccessJWTClaim) Valid() error {
	now := time.Now()
	if !c.VerifyExpiresAt(now, true) {
		return ErrAccessTokenExpired
	}

	if c.TokenType != constant.AccessTokenType {
		return ErrInvalidTokenType
	}

	return nil
}

func (c RefreshJWTClaim) Valid() error {
	now := time.Now()
	if !c.VerifyExpiresAt(now, true) {
		return ErrRefreshTokenExpired
	}

	if c.TokenType != constant.RefreshTokenType {
		return ErrInvalidTokenType
	}

	return nil
}

func (c StepUpJWTClaim) Valid() error {
	now := time.Now()
	if !c.VerifyExpiresAt(now, true) {
		return ErrStepUpTokenExpired
	}

	if c.TokenType != constant.StepUpTokenType {
		return ErrInvalidTokenType
	}

	return nil
}

func GenerateAccessToken(payload SignAccessTokenPayload, config dependency.Config) (*string, error) {
	expiresAt := time.Now().Add(time.Minute * time.Duration(config.Jwt.AccessTokenExpiration))
	now := time.Now()

	registeredClaims := jwt.RegisteredClaims{
		Issuer:    config.App.AppName,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	claims := AccessJWTClaim{
		RegisteredClaims: registeredClaims,
		UserId:           payload.UserID,
		IsSeller:         payload.IsSeller,
		TokenType:        constant.AccessTokenType,
	}

	accessToken := jwt.NewWithClaims(constant.JWTSigningMethod, claims)
	t, err := accessToken.SignedString([]byte(config.Jwt.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func GenerateRefreshToken(config dependency.Config) (*string, error) {
	expiresAt := time.Now().Add(time.Minute * time.Duration(config.Jwt.RefreshTokenExpiration))
	now := time.Now()

	registeredClaims := jwt.RegisteredClaims{
		Issuer:    config.App.AppName,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	claims := RefreshJWTClaim{
		RegisteredClaims: registeredClaims,
		TokenType:        constant.RefreshTokenType,
	}

	refreshToken := jwt.NewWithClaims(constant.JWTSigningMethod, claims)

	t, err := refreshToken.SignedString([]byte(config.Jwt.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func SignStepUpToken(config dependency.Config) (*string, error) {
	now := time.Now()
	duration := time.Minute * time.Duration(config.Jwt.StepUpTokenExpiration)
	expiresAt := now.Add(duration)

	registeredClaims := jwt.RegisteredClaims{
		Issuer:    config.App.AppName,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	claims := StepUpJWTClaim{
		RegisteredClaims: registeredClaims,
		TokenType:        constant.StepUpTokenType,
	}

	stepUpToken := jwt.NewWithClaims(constant.JWTSigningMethod, claims)
	secretBytes := []byte(config.Jwt.JWTSecret)
	t, err := stepUpToken.SignedString(secretBytes)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func ValidateAccessToken(generateToken string, config dependency.Config) (*jwt.Token, error) {
	computeFunction := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(config.Jwt.JWTSecret), nil
	}

	token, err := jwt.ParseWithClaims(generateToken, new(AccessJWTClaim), computeFunction)
	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			if e, ok := e.Inner.(*CustomError); ok {
				return nil, e
			}

			return nil, err
		}
	}

	return token, nil
}

func ValidateRefreshToken(refreshToken string, config dependency.Config) (*jwt.Token, error) {
	var computeFunction jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(config.Jwt.JWTSecret), nil
	}

	claim := new(RefreshJWTClaim)
	token, err := jwt.ParseWithClaims(refreshToken, claim, computeFunction)
	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			if e, ok := e.Inner.(*CustomError); ok {
				return nil, e
			}

			return nil, err
		}
	}

	return token, nil
}

func ValidateStepUpToken(stepUpToken string, config dependency.Config) (*jwt.Token, error) {
	var computeFunction jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		secretBytes := []byte(config.Jwt.JWTSecret)
		return secretBytes, nil
	}

	claim := new(StepUpJWTClaim)
	token, err := jwt.ParseWithClaims(stepUpToken, claim, computeFunction)
	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			if e, ok := e.Inner.(*CustomError); ok {
				return nil, e
			}

			return nil, err
		}
	}

	return token, nil
}

func ParseAccessTokenClaim(accessToken string, config dependency.Config) (*AccessJWTClaim, error) {
	token, _ := ValidateAccessToken(accessToken, config)
	if t, ok := token.Claims.(*AccessJWTClaim); ok {
		return t, nil
	}
	return nil, ErrInvalidToken
}
