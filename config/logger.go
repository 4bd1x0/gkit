package config

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerInstance struct {
	logger             *zap.Logger
	sugaredLogger      *zap.SugaredLogger
	errorLogger        *zap.Logger
	sugaredErrorLogger *zap.SugaredLogger
}

var (
	loggerInstance *LoggerInstance
	loggerOnce     sync.Once
)

func NewLogger() *LoggerInstance {
	loggerOnce.Do(func() {
		loggerInstance = createLogger()
	})
	return loggerInstance
}

func createLogger() *LoggerInstance {
	cfg := &config.Logger

	l := &LoggerInstance{}

	enableConsole := cfg.Console
	minLevel := parseLogLevel(cfg.Level)

	// Get the path from config file
	logsDir := GetLogsDir()

	// File encoder: no color
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	fileEncoderConfig.ConsoleSeparator = " | "
	fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig)

	// Console encoder: with color
	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // With color
	consoleEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	consoleEncoderConfig.ConsoleSeparator = " | "
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	l.createInfoLogger(minLevel, logsDir, enableConsole, fileEncoder, consoleEncoder)
	l.createErrorLogger(minLevel, logsDir, enableConsole, fileEncoder, consoleEncoder)

	return l
}

func (l *LoggerInstance) createInfoLogger(minLevel zapcore.Level, logsDir string, enableConsole bool, fileEncoder, consoleEncoder zapcore.Encoder) {
	if minLevel > zapcore.InfoLevel {
		l.logger = zap.NewNop()
		l.sugaredLogger = l.logger.Sugar()
		return
	}

	filename := filepath.Join(logsDir, "info.log")
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    200,
		MaxBackups: 10,
		MaxAge:     28,
		Compress:   true,
	})
	zc := zapcore.NewCore(fileEncoder, writeSyncer, minLevel)
	if enableConsole {
		zc = zapcore.NewTee(
			zc,
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), minLevel),
		)
	}
	l.logger = zap.New(zc, zap.AddCaller(), zap.AddCallerSkip(1))
	l.sugaredLogger = l.logger.Sugar()
}

func (l *LoggerInstance) createErrorLogger(minLevel zapcore.Level, logsDir string, enableConsole bool, fileEncoder, consoleEncoder zapcore.Encoder) {
	if minLevel > zapcore.ErrorLevel {
		l.errorLogger = zap.NewNop()
		l.sugaredErrorLogger = l.errorLogger.Sugar()
		return
	}

	filename := filepath.Join(logsDir, "error.log")
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    200,
		MaxBackups: 10,
		MaxAge:     28,
		Compress:   true,
	})
	zc := zapcore.NewCore(fileEncoder, writeSyncer, zap.ErrorLevel)
	if enableConsole {
		zc = zapcore.NewTee(
			zc,
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.ErrorLevel),
		)
	}
	l.errorLogger = zap.New(zc,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel), // Log stack trace at error level
	)
	l.sugaredErrorLogger = l.errorLogger.Sugar()
}

func (l *LoggerInstance) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *LoggerInstance) Debugf(template string, args ...interface{}) {
	l.sugaredLogger.Debugf(template, args...)
}

func (l *LoggerInstance) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *LoggerInstance) Infof(template string, args ...interface{}) {
	l.sugaredLogger.Infof(template, args...)
}

func (l *LoggerInstance) Error(msg string, fields ...zap.Field) {
	l.errorLogger.Error(msg, fields...)
}

func (l *LoggerInstance) Errorf(template string, args ...interface{}) {
	l.sugaredErrorLogger.Errorf(template, args...)
}

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

const (
	EnvLogDir = "GO_LOGS_DIR"
)

// GetLogsDir gets the log directory path
// Priority: 1. Environment variable GO_LOG_DIR > 2. logger.path in config file > 3. logs in current directory
func GetLogsDir() string {
	pathFromConfig := config.Logger.Path
	// 1. Priority: use log directory specified by environment variable
	if logsDir := os.Getenv(EnvLogDir); logsDir != "" {
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			panic("GO_LOGS_DIR is invalid or cannot create directory: " + err.Error())
		}
		return logsDir
	}

	// 2. Use logger.path from config file
	if logsDir := pathFromConfig; logsDir != "" {
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			panic("logger.path is invalid or cannot create directory: " + err.Error())
		}
		return logsDir
	}

	// 3. Logs in executable directory
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		logsDir := filepath.Join(exeDir, "logs")
		_, err := os.Stat(logsDir)
		if err == nil {
			return logsDir
		}
	}

	// 4. Logs in current directory
	logsDir := "./logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		panic("failed to create logs directory: " + err.Error())
	}
	return logsDir
}
