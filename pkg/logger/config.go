package logger

// Config godoc
type Config struct {
	Source string `mapstructure:"source"`
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}
