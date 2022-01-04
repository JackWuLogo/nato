package log

import (
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"micro-libs/utils/color"
	"os"
)

var (
	Logger       *logger.Helper
	disableColor = false
)

func init() {
	lvl, err := logger.GetLevel(os.Getenv("GAME_LOG_LEVEL"))
	if err != nil {
		lvl = logger.InfoLevel
	}

	Logger = logger.NewHelper(logger.NewLogger(
		logger.WithLevel(lvl),
		logger.WithCallerSkipCount(2),
	))
}

// 禁用颜色
func SetDisableColor() {
	disableColor = true
}

// 显示颜色
func SetEnableColor() {
	disableColor = false
}

// 获取当前日志级别
func GetLevel() logger.Level {
	return Logger.Options().Level
}

// 检查当前日志输出级别
func CheckLevel(level logger.Level) bool {
	return logger.V(level, Logger)
}

// 是否Trace模式
func IsTrace() bool {
	return CheckLevel(logger.TraceLevel)
}

// 是否开发模式
func IsDev() bool {
	return CheckLevel(logger.DebugLevel)
}

func Trace(tpl string, a ...interface{}) {
	Logger.Log(logger.TraceLevel, format(color.Trace, tpl, a...))
}

func Debug(tpl string, a ...interface{}) {
	Logger.Log(logger.DebugLevel, format(color.Debug, tpl, a...))
}

func Info(tpl string, a ...interface{}) {
	Logger.Log(logger.InfoLevel, format(color.Info, tpl, a...))
}

func Warn(tpl string, a ...interface{}) {
	Logger.Log(logger.WarnLevel, format(color.Warn, tpl, a...))
}

func Error(tpl string, a ...interface{}) {
	Logger.Log(logger.ErrorLevel, format(color.Error, tpl, a...))
}

func Fatal(tpl string, a ...interface{}) {
	Logger.Log(logger.FatalLevel, format(color.Fatal, tpl, a...))
}

func Comment(tpl string, a ...interface{}) {
	Logger.Log(logger.DebugLevel, format(color.Comment, tpl, a...))
}

func Success(tpl string, a ...interface{}) {
	Logger.Log(logger.InfoLevel, format(color.Success, tpl, a...))
}

func Question(tpl string, a ...interface{}) {
	Logger.Log(logger.WarnLevel, format(color.Question, tpl, a...))
}

// 格式化字符串
func format(theme *color.Theme, tpl string, a ...interface{}) string {
	if disableColor || theme == nil {
		return fmt.Sprintf(tpl, a...)
	}
	return theme.Text(tpl, a...)
}
