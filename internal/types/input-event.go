package types

type InputEvent struct {
	Data map[string]interface{} `mapstructure:"data" json:"data"`
}
