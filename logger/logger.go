package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	logger *slog.Logger
	level  slog.LevelVar
}

type LoggerOption func(*Logger)

// WithGroup 添加组
func WithGroup(name string) LoggerOption {
	return func(l *Logger) {
		l.logger = l.logger.WithGroup(name)
	}
}

type myHandler struct {
	w      io.Writer
	opts   slog.HandlerOptions
	groups []string
}

func NewGroupPrefixHandler(w io.Writer, opts *slog.HandlerOptions) *myHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &myHandler{
		w:    w,
		opts: *opts,
	}
}

func (h *myHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *myHandler) Handle(ctx context.Context, r slog.Record) error {
	// 构建组名字符串
	groupStr := ""
	if len(h.groups) > 0 {
		groupStr = strings.Join(h.groups, ".") + " "
	}

	// 格式化时间
	timeStr := r.Time.Format("2006/01/02 15:04:05")

	// 格式化级别
	levelStr := r.Level.String()
	if len(levelStr) > 0 {
		levelStr = strings.ToUpper(levelStr[:1]) + levelStr[1:]
	}

	// 根据日志级别添加颜色
	colorStart := ""
	colorEnd := ""
	switch r.Level {
	case slog.LevelDebug:
		colorStart = "\033[36m" // 青色
		colorEnd = "\033[0m"
	case slog.LevelInfo:
		colorStart = "\033[32m" // 绿色
		colorEnd = "\033[0m"
	case slog.LevelWarn:
		colorStart = "\033[33m" // 黄色
		colorEnd = "\033[0m"
	case slog.LevelError:
		colorStart = "\033[31m" // 红色
		colorEnd = "\033[0m"
	}

	// 应用相同颜色到组名和级别
	if groupStr != "" {
		groupStr = colorStart + groupStr + colorEnd
	}

	// 收集所有属性
	var attrs []string
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value))
		return true
	})

	// 构建输出行
	output := fmt.Sprintf("%s%s %s%s%s %s",
		groupStr, timeStr, colorStart, levelStr, colorEnd, r.Message)

	if len(attrs) > 0 {
		output += " " + strings.Join(attrs, " ")
	}
	output += "\n"

	// 写入输出
	_, err := h.w.Write([]byte(output))
	return err
}
func (h *myHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// 创建新 Handler 副本
	newHandler := &myHandler{
		w:      h.w,
		opts:   h.opts,
		groups: make([]string, len(h.groups)),
	}
	copy(newHandler.groups, h.groups)
	return newHandler
}

func (h *myHandler) WithGroup(name string) slog.Handler {
	// 创建新 Handler 副本
	newHandler := &myHandler{
		w:      h.w,
		opts:   h.opts,
		groups: make([]string, len(h.groups), len(h.groups)+1),
	}
	copy(newHandler.groups, h.groups)
	newHandler.groups = append(newHandler.groups, name)
	return newHandler
}

// NewLogger 创建一个新的 Logger 实例
func NewLogger(opts ...LoggerOption) *Logger {
	Logger := &Logger{
		level: slog.LevelVar{},
	}
	Logger.SetLevel(slog.LevelInfo)

	handler := NewGroupPrefixHandler(os.Stdout, &slog.HandlerOptions{
		Level: &Logger.level,
	})
	Logger.logger = slog.New(handler)

	for _, opt := range opts {
		opt(Logger)
	}
	return Logger
}
func (l *Logger) SetLevel(level slog.Level) *Logger {
	l.level.Set(level)
	return l
}
func (l *Logger) WithGroup(name string) *Logger {
	newlogger := &Logger{
		logger: l.logger.WithGroup(name),
		level:  slog.LevelVar{},
	}
	newlogger.SetLevel(l.level.Level())
	return newlogger
}

// Debug 记录 Debug 级别的日志
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info 记录 Info 级别的日志
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn 记录 Warn 级别的日志
func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error 记录 Error 级别的日志
func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
