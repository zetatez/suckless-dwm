package xlog

import (
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
	Compress   bool   `mapstructure:"compress"`
	Format     string `mapstructure:"format"`
	Console    bool   `mapstructure:"console"`
}

func NewLogger(cfg LogConfig) (*logrus.Logger, error) {
	logger := logrus.New()

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	if cfg.Format == "json" {
		logger.SetFormatter(
			&logrus.JSONFormatter{TimestampFormat: time.RFC3339},
		)
	} else {
		logger.SetFormatter(
			&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true},
		)
	}

	var outputs []io.Writer

	if cfg.Console {
		outputs = append(outputs, os.Stdout)
	}

	if cfg.Filename != "" {
		lj := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}
		outputs = append(outputs, lj)
	}

	if len(outputs) == 0 {
		outputs = append(outputs, os.Stdout)
	}

	logger.SetOutput(io.MultiWriter(outputs...))
	return logger, nil
}
