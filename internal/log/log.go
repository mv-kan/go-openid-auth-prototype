package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger = newZapLogger()
)

func Debug(msg string) {
	logger.Debug(msg)
}

func Info(msg string) {
	logger.Info(msg)
}

func Error(msg string) {
	logger.Error(msg)
}

func Fatal(msg string) {
	logger.Fatal(msg)
}

func newZapLogger() *zap.Logger {
	zapConfig := zap.NewProductionEncoderConfig()
	zapConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(zapConfig)

	f, err := os.OpenFile("./main-api.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fileEncoder := zapcore.NewJSONEncoder(zapConfig)

	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), defaultLogLevel),
	)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.InfoLevel))

	return zapLogger
}
