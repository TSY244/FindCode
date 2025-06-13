package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	defaultLogger *zap.SugaredLogger
)

func init() {
	defaultLoggerInit()
}

func defaultLoggerInit() {
	defaultLogger = zap.NewExample().Sugar()
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
}

func Infoln(args ...interface{}) {
	defaultLogger.Infoln(args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	defaultLogger.Debugf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	defaultLogger.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	defaultLogger.Panicf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	defaultLogger.Warnf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	defaultLogger.Fatalf(template, args...)
}

func GetDefaultLogger() *zap.SugaredLogger {
	return defaultLogger
}

func SetDefaultLogger(logger *zap.SugaredLogger) {
	defaultLogger = logger
}

func NewLogger(config *LogConf) *zap.SugaredLogger {
	if config == nil {
		return defaultLogger
	}

	// 定义日志输出格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 启用颜色
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // 使用 ISO8601 时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 使用 ConsoleEncoder 代替 JSONEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.FileName,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	})

	var level zapcore.Level
	switch config.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "dpanic":
		level = zapcore.DPanicLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	// 文件日志核心
	fileCore := zapcore.NewCore(encoder, fileWriteSyncer, zap.NewAtomicLevelAt(level))

	// 如果需要同时输出到控制台
	if config.PrintToConsole {
		consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(level))
		core := zapcore.NewTee(fileCore, consoleCore)
		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		return logger.Sugar()
	} else {
		logger := zap.New(fileCore, zap.AddCaller(), zap.AddCallerSkip(1))
		return logger.Sugar()
	}
}
