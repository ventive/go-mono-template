package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger godoc
type Logger struct {
	section string
	action  string
	meta    map[string]interface{}
}

// New godoc
func New(section, action string) *Logger {
	return &Logger{
		section: section,
		action:  action,
		meta:    map[string]interface{}{},
	}
}

// AddMeta godoc
func (l *Logger) AddMeta(name string, meta interface{}) {
	l.meta[name] = meta
}

// Zerolog adds section and action to the given zerolog event and returns it
func (l *Logger) Zerolog(event *zerolog.Event) *zerolog.Event {
	event = event.Str("section", l.section).Str("action", l.action)
	if len(l.meta) > 0 {
		event = event.Interface("meta", l.meta)
	}

	return event
}

// Info godoc
func (l *Logger) Info(msg string) {
	l.Zerolog(log.Info()).Msg(msg)
}

// Debug godoc
func (l *Logger) Debug(msg string) {
	l.Zerolog(log.Debug()).Msg(msg)
}

// DebugWithExtra godoc
func (l *Logger) DebugWithExtra(msg string, extra map[string]interface{}) {
	evt := log.Debug()
	for k, v := range extra {
		evt.Interface(k, v)
	}

	l.Zerolog(evt).Msg(msg)
}

// Error godoc
func (l *Logger) Error(msg string, errs ...error) {
	evt := log.Error()
	for _, err := range errs {
		evt = evt.Err(err)
	}

	l.Zerolog(evt).Msg(msg)
}

// ErrorWithExtra godoc
func (l *Logger) ErrorWithExtra(msg string, extra map[string]interface{}, errs ...error) {
	evt := log.Error()

	for k, v := range extra {
		evt.Interface(k, v)
	}

	for _, err := range errs {
		evt = evt.Err(err)
	}

	l.Zerolog(evt).Msg(msg)
}

// Fatal godoc
func (l *Logger) Fatal(msg string, errs ...error) {
	evt := log.Fatal()
	for _, err := range errs {
		evt = evt.Err(err)
	}

	l.Zerolog(evt).Msg(msg)
}

// Warn godoc
func (l *Logger) Warn(msg string) {
	l.Zerolog(log.Warn()).Msg(msg)
}

// Init godoc
func Init(cfg Config) {
	if strings.ToLower(cfg.Format) == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.SetGlobalLevel(stringLevelToZerologLevel(cfg.Level))

	log.Info().Str("log_level", cfg.Level).
		Str("log_format", cfg.Format).
		Msg("Logger level")
}

func stringLevelToZerologLevel(level string) zerolog.Level {
	level = strings.ToLower(level)
	levelMap := map[string]zerolog.Level{
		"trace": zerolog.TraceLevel,
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
	}

	if zerologLevel, ok := levelMap[level]; ok {
		return zerologLevel
	}

	return zerolog.InfoLevel
}
