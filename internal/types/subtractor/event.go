package subtractor

type SubtractEvent struct {
	A float64 `mapstructure:"a" json:"a" validate:"number,required"`
	B float64 `mapstructure:"b" json:"b" validate:"number,required"`
}
