package shared

import (
	"math"

	"github.com/google/uuid"
	"github.com/jaevor/go-nanoid"
	"github.com/lil-oren/rest/internal/constant"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateNanoID() string {
	canonicID, err := nanoid.CustomASCII(constant.VerifCodeAlphaNum, 6)
	if err != nil {
		return ""
	}
	otp := canonicID()
	return otp
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
