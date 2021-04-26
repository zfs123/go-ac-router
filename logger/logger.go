package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger      *zap.Logger
	atomicLevel zap.AtomicLevel
)

// log rotation settings
func getLogWriter(filename string, maxSize, maxBackups, maxAge int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// json serialized log format
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// set log level
func setLogLevel(logLevel string) error {
	atomicLevel = zap.NewAtomicLevel()
	var level = zapcore.DebugLevel
	err := level.Set(logLevel)
	if err != nil {
		return err
	}
	atomicLevel.SetLevel(level)
	return nil
}

// Initialize logger according to configuration
func InitLogger(filename string, maxSize, maxBackups, maxAge int, compress bool, logLevel string) error {
	writeSyncer := getLogWriter(filename, maxSize, maxBackups, maxAge, compress)
	encoder := getEncoder()
	setLogLevel(logLevel)

	core := zapcore.NewCore(encoder, writeSyncer, atomicLevel)
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(1)

	logger = zap.New(core, caller, zap.AddStacktrace(zap.PanicLevel), callerSkip)
	return nil
}
