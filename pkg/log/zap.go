package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2" // Include this for lumberjack
)

// Config is the configuration for the logger.
type Config struct {
	// StdOutLogLevel is the log level for the logger.
	StdOutLogLevel string
	// FileOutLogLevel is the log level for the file logger.
	FileOutLogLevel string
	// WriteTo is the output file for the logger. If empty, logs will be written to stderr.
	WriteTo string
	// MaxSize is the maximum size in megabytes before log is rotated.
	MaxSize int
	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int
	// MaxAge is the maximum number of days to retain an old log file.
	MaxAge int
	// Compress determines if the rotated log files should be compressed.
	Compress bool
}

// NewDefaultConfig creates a default configuration for the logger.
func NewDefaultConfig() Config {
	return Config{
		StdOutLogLevel:  "info",
		FileOutLogLevel: "info",
		WriteTo:         "",
		MaxSize:         100, // 100MB
		MaxBackups:      2,
		MaxAge:          3, // 1 day
		Compress:        false,
	}
}

func NewLogger(config Config) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// Setup the primary output to always include os.Stderr
	stderrSyncer := zapcore.Lock(os.Stderr)

	var fileSyncer zapcore.WriteSyncer
	if config.WriteTo != "" && config.WriteTo != "stderr" {
		// Configure lumberjack for logging to a file
		lumberjackLogger := &lumberjack.Logger{
			Filename:   config.WriteTo,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		fileSyncer = zapcore.AddSync(lumberjackLogger)
	}

	// Use zapcore.NewTee to write to both stderr and the file (if configured)
	var fileCore zapcore.Core
	if fileSyncer != nil {
		logLevel := zapcore.InfoLevel // Default log level
		if err := logLevel.Set(config.FileOutLogLevel); err != nil {
			fmt.Fprintf(os.Stderr, "failed to set log level on file logging: %v\nfalling back to info", err)
			logLevel = zapcore.InfoLevel // Fallback to info if setting fails
		}

		fileCore = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			fileSyncer,
			logLevel,
		)
	}

	logLevel := zapcore.InfoLevel // Default log level
	if err := logLevel.Set(config.StdOutLogLevel); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set log level on std out: %v\nfalling back to info", err)
		logLevel = zapcore.InfoLevel // Fallback to info if setting fails
	}

	stdCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		stderrSyncer,
		logLevel,
	)

	var core zapcore.Core
	if fileCore != nil {
		core = zapcore.NewTee(stdCore, fileCore)
	} else {
		core = stdCore
	}

	return zap.New(
		core,
		zap.AddCaller(),
		zap.Fields(zapcore.Field{Key: "pid", Type: zapcore.Int64Type, Integer: int64(os.Getpid())}),
	)
}
