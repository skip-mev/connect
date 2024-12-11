package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2" // Include this for lumberjack
)

// Config is the configuration for the logger.
type Config struct {
	// StdOutLogLevel is the log level for the standard out logger.
	StdOutLogLevel string
	// FileOutLogLevel is the log level for the file logger.
	FileOutLogLevel string
	// DisableRotating disables log rotation.
	DisableRotating bool
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
	// LogSamplePeriod is the duration in which we de-dupe identical log messages.
	LogSamplePeriod time.Duration
}

// NewDefaultConfig creates a default configuration for the logger.
func NewDefaultConfig() Config {
	return Config{
		StdOutLogLevel:  "info",
		FileOutLogLevel: "info",
		DisableRotating: false,
		WriteTo:         "sidecar.log",
		MaxSize:         1, // 100MB
		MaxBackups:      1,
		MaxAge:          3, // 3 days
		Compress:        false,
		LogSamplePeriod: 10 * time.Second,
	}
}

func NewLogger(config Config) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var fileCore zapcore.Core
	if config.WriteTo != "" && !config.DisableRotating {
		// Configure lumberjack for logging to a file
		lumberjackLogger := &lumberjack.Logger{
			Filename:   config.WriteTo,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		fileSyncer := zapcore.AddSync(lumberjackLogger)

		logLevel := zapcore.InfoLevel
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

	// Setup the primary output to always include os.Stderr.
	logLevel := zapcore.InfoLevel
	if err := logLevel.Set(config.StdOutLogLevel); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set log level on std out: %v\nfalling back to info", err)
		logLevel = zapcore.InfoLevel // Fallback to info if setting fails
	}

	// Setup the primary output to always include os.Stderr
	stdCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stderr),
		logLevel,
	)

	// Use zapcore.NewTee to write to both stderr and the file (if configured)
	var core zapcore.Core
	if fileCore != nil {
		core = zapcore.NewTee(stdCore, fileCore)
	} else {
		core = stdCore
	}
	if strings.ToUpper(config.StdOutLogLevel) != zap.DebugLevel.CapitalString() && strings.ToUpper(config.FileOutLogLevel) != zap.DebugLevel.CapitalString() {
		// If we're not in debug log level anywhere, filter any logs which have non-unique messages within a 10-second period
		core = zapcore.NewSamplerWithOptions(core, config.LogSamplePeriod, 1, 0)
	}

	return zap.New(
		core,
		zap.AddCaller(),
		zap.Fields(zapcore.Field{Key: "pid", Type: zapcore.Int64Type, Integer: int64(os.Getpid())}),
	)
}
