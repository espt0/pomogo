package validator

import (
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

var (
	instance *validator.Validate
	once     sync.Once
)

type CustomValidator struct{}

func NewEchoValidator() *CustomValidator {
	return &CustomValidator{}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := get().Struct(i); err != nil { // usa o singleton
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func get() *validator.Validate {

	once.Do(func() {
		instance = validator.New(validator.WithRequiredStructEnabled())

		instance.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	})

	return instance
}
