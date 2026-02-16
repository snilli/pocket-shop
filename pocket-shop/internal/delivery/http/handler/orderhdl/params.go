package orderhdl

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type GetOrderParams struct {
	ID string `validate:"required,uuid"`
}
