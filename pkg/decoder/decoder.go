package decoder

import (
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

var validate = validator.New()

// Decode takes a map and translate it to a struct.
// Validates the struct using go-playground/validator pkg.
//
// event must be a pointer to a struct with "validate" tags attached
// "validate" tags are optional. If not provided, struct will be always valid.
func Decode(data map[string]interface{}, event interface{}) error {
	err := mapstructure.Decode(data, event)
	if err != nil {
		return err
	}
	return validate.Struct(event)
}
