package internallogger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(level string) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.TimeKey = "time"

	var logLevel zapcore.Level
	var err error

	if level != "" {
		logLevel, err = zapcore.ParseLevel(level)
		if err != nil {
			log.Println("unable to set level") // TODO
			logLevel = zap.InfoLevel
			// return err
		}
	}

	logg := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(logLevel),
	))

	defer logg.Sync() //nolint

	logg.Info("start logging")

	// logg.Info(GetVersion())
	return logg
}
