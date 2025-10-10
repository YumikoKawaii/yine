package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	prefix        string
	sugaredLogger *zap.SugaredLogger
	writer        *lumberjack.Logger
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(l.prefix+format, args...)
}

func (l *zapLogger) Debug(msg string) {
	l.sugaredLogger.Debug(l.prefix + msg)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(l.prefix+format, args...)
}

func (l *zapLogger) Info(msg string) {
	l.sugaredLogger.Info(l.prefix + msg)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(l.prefix+format, args...)
}

func (l *zapLogger) Warn(msg string) {
	l.sugaredLogger.Warn(l.prefix + msg)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(l.prefix+format, args...)
}

func (l *zapLogger) Error(msg string) {
	l.sugaredLogger.Error(l.prefix + msg)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(l.prefix+format, args...)
}

func (l *zapLogger) Fatal(msg string) {
	l.sugaredLogger.Fatal(l.prefix + msg)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(l.prefix+format, args...)
}

func (l *zapLogger) Panic(msg string) {
	l.sugaredLogger.Panic(l.prefix + msg)
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{
		prefix:        l.prefix,
		sugaredLogger: newLogger,
		writer:        l.writer,
	}
}

func (l *zapLogger) WithPrefix(prefix string) Logger {
	newLogger := l.sugaredLogger.With()
	return &zapLogger{
		prefix:        l.prefix + prefix,
		sugaredLogger: newLogger,
		writer:        l.writer,
	}
}

func (l *zapLogger) WithTraceId(ctx context.Context) Logger {
	return &zapLogger{
		prefix:        l.prefix,
		sugaredLogger: l.sugaredLogger,
		writer:        l.writer,
	}
}

func (l *zapLogger) GetDelegate() interface{} {
	return l.sugaredLogger
}

func (l *zapLogger) CloseWriter() error {
	if l.writer != nil {
		return l.writer.Close()
	}

	return nil
}

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case infoLvl:
		return zapcore.InfoLevel
	case warnLvl:
		return zapcore.WarnLevel
	case debugLvl:
		return zapcore.DebugLevel
	case errorLvl:
		return zapcore.ErrorLevel
	case fatalLvl:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func newZapLogger(config Configuration) (Logger, error) {
	var cores []zapcore.Core

	if config.EnableConsole {
		level := getZapLevel(config.ConsoleLevel)
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(config.ConsoleJSONFormat), writer, level)
		cores = append(cores, core)
	}

	var writer *lumberjack.Logger
	if config.EnableFile {
		level := getZapLevel(config.FileLevel)
		writer = &lumberjack.Logger{
			Filename: config.FileLocation,
			MaxSize:  100,
			Compress: true,
			MaxAge:   28,
		}
		writeSyncer := zapcore.AddSync(writer)
		core := zapcore.NewCore(getEncoder(config.FileJSONFormat), writeSyncer, level)
		cores = append(cores, core)
	}

	// AddCallerSkip skips 2 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	logger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	)

	return &zapLogger{
		sugaredLogger: logger.Sugar(),
		writer:        writer,
	}, nil
}

func GetZapLoggerDelegate(instance Logger) (*zap.SugaredLogger, error) {
	switch v := instance.GetDelegate().(type) {
	case *zap.SugaredLogger:
		return instance.GetDelegate().(*zap.SugaredLogger), nil
	default:
		return nil, fmt.Errorf("expected zap.SugaredLogger but got: %v", v)
	}
}
