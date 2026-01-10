package logger

import (
	"context"
	"ecommerce/gateway/config"
	"io"
	"os"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlog "github.com/cloudwego/hertz/pkg/common/hlog"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger 实现Hertz的FullLogger接口
type ZapLogger struct {
	logger *zap.Logger
	level  hertzlog.Level
}

// SetLevel 设置日志级别
func (z *ZapLogger) SetLevel(level hertzlog.Level) {
	z.level = level
}

// GetLevel 获取日志级别
func (z *ZapLogger) GetLevel() hertzlog.Level {
	return z.level
}

// IsLevelEnabled 检查级别是否启用
func (z *ZapLogger) IsLevelEnabled(level hertzlog.Level) bool {
	return level >= z.level
}

// Log 记录日志
func (z *ZapLogger) Log(level hertzlog.Level, v ...interface{}) {
	switch level {
	case hertzlog.LevelTrace, hertzlog.LevelDebug:
		z.logger.Debug("", zap.Any("msg", v))
	case hertzlog.LevelInfo, hertzlog.LevelNotice:
		z.logger.Info("", zap.Any("msg", v))
	case hertzlog.LevelWarn:
		z.logger.Warn("", zap.Any("msg", v))
	case hertzlog.LevelError:
		z.logger.Error("", zap.Any("msg", v))
	case hertzlog.LevelFatal:
		z.logger.Fatal("", zap.Any("msg", v))
	}
}

// Logf 格式化记录日志
func (z *ZapLogger) Logf(level hertzlog.Level, format string, v ...interface{}) {
	z.Log(level, v...)
}

// CtxLogf 带上下文的日志
func (z *ZapLogger) CtxLogf(ctx context.Context, level hertzlog.Level, format string, v ...interface{}) {
	z.Log(level, v...)
}

// Trace 记录跟踪日志
func (z *ZapLogger) Trace(v ...interface{}) {
	z.Log(hertzlog.LevelTrace, v...)
}

// Debug 记录调试日志
func (z *ZapLogger) Debug(v ...interface{}) {
	z.Log(hertzlog.LevelDebug, v...)
}

// Info 记录信息日志
func (z *ZapLogger) Info(v ...interface{}) {
	z.Log(hertzlog.LevelInfo, v...)
}

// Notice 记录通知日志
func (z *ZapLogger) Notice(v ...interface{}) {
	z.Log(hertzlog.LevelNotice, v...)
}

// Warn 记录警告日志
func (z *ZapLogger) Warn(v ...interface{}) {
	z.Log(hertzlog.LevelWarn, v...)
}

// Error 记录错误日志
func (z *ZapLogger) Error(v ...interface{}) {
	z.Log(hertzlog.LevelError, v...)
}

// Fatal 记录致命错误日志
func (z *ZapLogger) Fatal(v ...interface{}) {
	z.Log(hertzlog.LevelFatal, v...)
}

// Tracef 格式化跟踪日志
func (z *ZapLogger) Tracef(format string, v ...interface{}) {
	z.Logf(hertzlog.LevelTrace, format, v...)
}

// Debugf 格式化调试日志
func (z *ZapLogger) Debugf(format string, v ...interface{}) {
	z.Logf(hertzlog.LevelDebug, format, v...)
}

// Infof 格式化信息日志
func (z *ZapLogger) Infof(format string, v ...interface{}) {
	z.Logf(hertzlog.LevelInfo, format, v...)
}

// Noticef 格式化通知日志
func (z *ZapLogger) Noticef(format string, v ...interface{}) {
	z.Logf(hertzlog.LevelNotice, format, v...)
}

// Warnf 格式化警告日志
func (z *ZapLogger) Warnf(format string, v ...interface{}) {
	z.Logf(hertzlog.LevelWarn, format, v...)
}

// Errorf 格式化错误日志
func (z *ZapLogger) Errorf(format string, v ...interface{}) {
	z.Logf(hertzlog.LevelError, format, v...)
}

// Fatalf 格式化致命错误日志
func (z *ZapLogger) Fatalf(format string, v ...interface{}) {
	z.Logf(hertzlog.LevelFatal, format, v...)
}

// CtxTracef 带上下文的跟踪日志
func (z *ZapLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	z.CtxLogf(ctx, hertzlog.LevelTrace, format, v...)
}

// CtxDebugf 带上下文的调试日志
func (z *ZapLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	z.CtxLogf(ctx, hertzlog.LevelDebug, format, v...)
}

// CtxInfof 带上下文的信息日志
func (z *ZapLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	z.CtxLogf(ctx, hertzlog.LevelInfo, format, v...)
}

// CtxNoticef 带上下文的通知日志
func (z *ZapLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	z.CtxLogf(ctx, hertzlog.LevelNotice, format, v...)
}

// CtxWarnf 带上下文的警告日志
func (z *ZapLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	z.CtxLogf(ctx, hertzlog.LevelWarn, format, v...)
}

// CtxErrorf 带上下文的错误日志
func (z *ZapLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	z.CtxLogf(ctx, hertzlog.LevelError, format, v...)
}

// CtxFatalf 带上下文的致命错误日志
func (z *ZapLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	z.CtxLogf(ctx, hertzlog.LevelFatal, format, v...)
}

// InitLogger 初始化日志系统
func InitLogger(cfg *config.Config) {
	// 创建日志目录
	if err := os.MkdirAll(cfg.Log.Path, 0755); err != nil {
		// 如果无法创建目录，使用标准错误输出
		hlog.SetOutput(os.Stderr)
		hlog.Errorf("创建日志目录失败: %v", err)
	}

	// 设置Hertz的日志级别
	setHertzLogLevel(cfg.Log.Level)
	initSimpleLogger(cfg)
}

// initSimpleLogger 使用简单的hlog配置
func initSimpleLogger(cfg *config.Config) {
	var output io.Writer

	// 配置日志输出
	switch cfg.Log.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	case "file":
		logFile := filepath.Join(cfg.Log.Path, "gateway.log")
		output = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    100, // MB
			MaxBackups: 10,
			MaxAge:     30, // days
			Compress:   true,
		}
	default:
		output = os.Stdout
	}

	// 设置日志输出
	hlog.SetOutput(output)

	// 记录初始化成功
	hlog.Infof("✅ 日志系统初始化完成 - 级别: %s, 输出: %s", cfg.Log.Level, cfg.Log.Output)
}

// setHertzLogLevel 设置Hertz的日志级别
func setHertzLogLevel(level string) {
	switch level {
	case "debug":
		hlog.SetLevel(hlog.LevelDebug)
	case "info":
		hlog.SetLevel(hlog.LevelInfo)
	case "warn":
		hlog.SetLevel(hlog.LevelWarn)
	case "error":
		hlog.SetLevel(hlog.LevelError)
	default:
		hlog.SetLevel(hlog.LevelInfo)
	}
}
