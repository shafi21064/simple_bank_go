package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/shafi21064/simplebank/util"
)

var validateCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupporpedCurrency(currency)
	}
	return false
}
