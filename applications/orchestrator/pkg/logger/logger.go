package logger

import (
	"context"
	"errors"
	"sync"
)

// A global variable so that log functions can be directly accessed
var log Logger = DefaultLogger()

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// LoggerBackend reprents the int enum for backend of logger (exists for legacy compat reason)
type LoggerBackend string

type LogLvl string

const (
	// Debug has verbose message
	DebugLvl LogLvl = "debug"
	//Info is default log level
	InfoLvl LogLvl = "info"
	// Warn is for logging messages about possible issues
	WarnLvl LogLvl = "warn"
	// Error is for logging errors
	ErrorLvl LogLvl = "error"
	// Fatal is for logging fatal messages. The sytem shutsdown after logging the message.
	FatalLvl LogLvl = "fatal"
)

const (
	// Debug has verbose message
	debugLvl = "debug"
	//Info is default log level
	infoLvl = "info"
	// Warn is for logging messages about possible issues
	warnLvl = "warn"
	// Error is for logging errors
	errorLvl = "error"
	// Fatal is for logging fatal messages. The sytem shutsdown after logging the message.
	fatalLvl = "fatal"
)

const (
	// LoggerBackendZap logging using Uber's zap backend
	LoggerBackendZap LoggerBackend = "zap"
	// LoggerBackendLogrus logging using logrus backend
	LoggerBackendLogrus LoggerBackend = "logrus"
)

var (
	errInvalidLoggerInstance = errors.New("invalid logger instance")

	once sync.Once
)

// Logger is our contract for the logger
type Logger interface {
	Debug(msg string)
	Debugf(format string, args ...interface{})

	Info(msg string)
	Infof(format string, args ...interface{})

	Warn(msg string)
	Warnf(format string, args ...interface{})

	Error(msg string)
	Errorf(format string, args ...interface{})

	Fatal(msg string)
	Fatalf(format string, args ...interface{})

	Panic(msg string)
	Panicf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger
	WithPrefix(prefix string) Logger
	WithTraceId(ctx context.Context) Logger

	GetDelegate() interface{}

	CloseWriter() error
}

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration struct {
	Backend           LoggerBackend `json:"backend" mapstructure:"backend" name:"log-backend" help:"Logger backend" env:"BACKEND" default:"zap" enum:"zap, logrus"`
	EnableConsole     bool          `json:"enable_console" mapstructure:"enable_console" name:"log-enable-console" help:"Enable log console" env:"ENABLE_CONSOLE" default:"true"`
	ConsoleJSONFormat bool          `json:"console_json_format" mapstructure:"console_json_format" name:"log-console-json-format" help:"Console to json format" env:"CONSOLE_JSON_FORMAT" default:"false"`
	ConsoleLevel      string        `json:"console_level" mapstructure:"console_level" name:"log-console-level" help:"Console log level" env:"CONSOLE_LEVEL" default:"info" enum:"debug, info, warn, error, fatal, panic"`
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}

func DefaultConfig() Configuration {
	return Configuration{
		Backend:           LoggerBackendZap,
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      "info",
		EnableFile:        false,
		FileJSONFormat:    false,
	}
}

// DefaultLogger creates default logger, which uses zap sugarlogger and outputs to console
func DefaultLogger() Logger {
	cfg := DefaultConfig()
	logger, _ := newZapLogger(cfg)
	return logger
}

// InitLogger returns an instance of logger
func InitLogger(conf Configuration) (Logger, error) {
	var err error
	once.Do(func() {
		switch conf.Backend {
		case LoggerBackendZap, LoggerBackendLogrus:
			log, err = NewLogger(conf)

		default:
			err = errInvalidLoggerInstance
		}
	})
	return log, err
}

func NewLogger(conf Configuration) (Logger, error) {
	return newZapLogger(conf)
}

func Debug(msg string) {
	log.Debugf(msg)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(msg string) {
	log.Infof(msg)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(msg string) {
	log.Warnf(msg)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(msg string) {
	log.Errorf(msg)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(msg string) {
	log.Fatalf(msg)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panic(msg string) {
	log.Panicf(msg)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

func WithPrefix(prefix string) Logger {
	return log.WithPrefix(prefix)
}

func WithTraceId(ctx context.Context) Logger {
	return log.WithTraceId(ctx)
}

func Get() Logger {
	return log
}

func GetDelegate() interface{} {
	return log.GetDelegate()
}
