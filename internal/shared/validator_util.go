package shared

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidatorUseJSONName(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
}
