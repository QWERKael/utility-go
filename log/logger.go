package log

import (
	"fmt"
	"github.com/QWERKael/utility-go/path"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var SugarLogger *zap.SugaredLogger

type EncoderFormat uint8

const (
	ConsoleEncoder EncoderFormat = 0
	JSONEncoder    EncoderFormat = 1
)

func InitLogger(logPath, LogLevel string) *zap.SugaredLogger {
	now := time.Now().Format("2006_01_02_15")
	if logPath != "" {
		logPath = fmt.Sprintf("%s.%s", logPath, now)
	}
	sugarLogger, err := NewLogger(ConsoleEncoder, logPath, ConvertLogLevel(LogLevel))
	if err != nil {
		panic(err.Error())
	}
	SugarLogger = sugarLogger
	sugarLogger.Debug("日志记录开始...")
	return sugarLogger
}

func NewLogger(encoderFormat EncoderFormat, logFilePath string, level zapcore.Level) (*zap.SugaredLogger, error) {
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	switch encoderFormat {
	case ConsoleEncoder:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case JSONEncoder:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var writeSyncer zapcore.WriteSyncer
	if logFilePath == "" {
		writeSyncer = zapcore.Lock(os.Stdout)
	} else {
		file, err := path.CreateOrOpenFileForAppendWrite(logFilePath)
		if err != nil {
			return nil, err
		}
		writeSyncer = zapcore.AddSync(file)
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger := logger.Sugar()
	return sugarLogger, nil
}

func ConvertLogLevel(strLevel string) zapcore.Level {
	switch strings.ToLower(strLevel) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	}
	return zapcore.InfoLevel
}
