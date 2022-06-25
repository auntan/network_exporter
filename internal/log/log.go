package log

import (
	"fmt"
	"github.com/auntan/network_exporter/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLoggerDefault() {
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic(fmt.Errorf("initialize logger error: %v", err))
	}
	zap.ReplaceGlobals(logger)
}

func InitLogger(conf *config.Config) error {
	var cfg zap.Config
	if conf.LogsEnv == "prod" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.MessageKey = "message"
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)

	return nil
}
