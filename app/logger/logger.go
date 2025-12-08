package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(service string, isDev bool) {
	var logger *zap.Logger
	var err error

	if isDev {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoderCfg.ConsoleSeparator = " "

		// Custom key names and ordering
		encoderCfg.TimeKey = "T"
		encoderCfg.LevelKey = "L"
		encoderCfg.NameKey = "N"
		encoderCfg.CallerKey = "" // Disable caller
		encoderCfg.MessageKey = "M"
		encoderCfg.StacktraceKey = "S"

		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		)
		logger = zap.New(core)
		logger = logger.Named(fmt.Sprintf("\x1b[35m[%s]\x1b[0m", service))
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
		logger = logger.With(zap.String("service", service))
	}

	Log = logger
}

func Info(msg string, fields ...zap.Field)  { Log.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { Log.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Log.Error(msg, fields...) }
func Debug(msg string, fields ...zap.Field) { Log.Debug(msg, fields...) }

type GinWriter struct{}

func (g GinWriter) Write(p []byte) (n int, err error) {
	Log.Info(string(p))
	return len(p), nil
}
