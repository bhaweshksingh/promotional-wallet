package reporters

import (
	"account/pkg/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

const permissions os.FileMode = 0744

func NewExternalLogFile(cfg config.LogFileConfig) *lumberjack.Logger {
	if err := os.MkdirAll(cfg.GetFileDir(), permissions); err != nil {
		return nil
	}

	return &lumberjack.Logger{
		Filename:   cfg.GetFilePath(),
		MaxSize:    cfg.GetFileMaxSizeInMb(),
		MaxBackups: cfg.GetFileMaxBackups(),
		MaxAge:     cfg.GetFileMaxAge(),
		LocalTime:  cfg.GetFileWithLocalTimeStamp(),
	}
}
