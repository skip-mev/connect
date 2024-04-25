package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config is the configuration for the logger.
type Config struct {
	// LogLevel is the log level to use. The default is "info".
	LogLevel string
	// WriteTo is the path to write logs to. The default is stderr.
	WriteTo string
}

// NewLogger creates a new logger with the given configuration. This logger
// is pre-configured to production settings. To change the settings, modify
// the Config struct.
func NewLogger(config Config) *zap.Logger {
	// EncodeTime is set to ISO8601TimeEncoder by default. This is a more human readable
	// format than the default EpochTimeEncoder.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// outputPaths is the list of paths to write logs to. By default, logs are
	// written to stderr.
	outputPaths := []string{"stderr"}
	if config.WriteTo != "" {
		outputPaths = append(outputPaths, config.WriteTo)
	}

	// initialFields are the fields that are added to every log message.
	initialFields := map[string]interface{}{
		"pid": os.Getpid(),
	}
	if config.WriteTo != "" {
		initialFields["writing_logs_to"] = config.WriteTo
	}

	var (
		logLevel zapcore.Level
		dev      = false
	)
	switch config.LogLevel {
	case "debug":
		logLevel = zapcore.DebugLevel
		dev = true
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	case "dpanic":
		logLevel = zapcore.DPanicLevel
	case "panic":
		logLevel = zapcore.PanicLevel
	case "fatal":
		logLevel = zapcore.FatalLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(logLevel),
		Development:       dev,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  outputPaths,
		InitialFields:     initialFields,
	}

	return zap.Must(zapConfig.Build())
}
