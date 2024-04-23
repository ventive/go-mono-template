package configparser

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testConfigFile      = "./.fixture/config.yml"
	defaultField1       = "default-field-1"
	defaultNestedField1 = "default-nested-field-1"
	envField1           = "env-field-1"
	envNestedField1     = "env-nested-field-1"
)

var defaultConfigValues = map[string]interface{}{
	"field1":        defaultField1,
	"nested.field1": defaultNestedField1,
}

type nested struct {
	Field1 string `mapstructure:"field1"`
}

type config struct {
	Field1 string `mapstructure:"field1"`
	Nested nested `mapstructure:"nested"`
}

func TestParse(t *testing.T) {
	t.Run("returnsErrWhenTargetNotPointer", func(t *testing.T) {
		cfg := config{}
		err := Parse(testConfigFile, cfg)
		assert.Equal(t, ErrTargetNotPointer, err)
	})

	t.Run("returnsErrWhenGivenFileDoesNotExist", func(t *testing.T) {
		cfg := config{}
		err := Parse(strconv.Itoa(int(time.Now().UnixNano())), &cfg)
		assert.Equal(t, ErrFileDoesNotExist, err)
	})

	t.Run("successWhenLoadingFromFile", func(t *testing.T) {
		cfg := config{}
		err := Parse(testConfigFile, &cfg)
		assert.Nil(t, err)

		assert.Equal(t, "fixture-field-1", cfg.Field1)
		assert.Equal(t, "fixture-nested-field-1", cfg.Nested.Field1)
	})

	t.Run("successWhenNoSourceVals", func(t *testing.T) {
		cfg := config{}
		err := Parse("", &cfg)
		assert.Nil(t, err)

		assert.Equal(t, "", cfg.Field1)
		assert.Equal(t, "", cfg.Nested.Field1)
	})

	t.Run("successWhenLoadingFromDefaults", func(t *testing.T) {
		cfg := config{}
		err := Parse("", &cfg, defaultConfigValues)
		assert.Nil(t, err)

		assert.Equal(t, defaultField1, cfg.Field1)
		assert.Equal(t, defaultNestedField1, cfg.Nested.Field1)
	})

	t.Run("successWhenLoadingFromEnv", func(t *testing.T) {
		_ = os.Setenv("FIELD1", envField1)
		_ = os.Setenv("NESTED_FIELD1", envNestedField1)

		cfg := config{}
		// should be able to get them without the defaultConfigValues but:
		// https://github.com/spf13/viper/issues/584#issuecomment-451554896
		err := Parse("", &cfg, defaultConfigValues)
		assert.Nil(t, err)

		assert.Equal(t, envField1, cfg.Field1)
		assert.Equal(t, envNestedField1, cfg.Nested.Field1)
	})
}
