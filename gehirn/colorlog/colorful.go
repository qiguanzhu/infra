package colorlog

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func (l *Logger) init() {
	l.debugLogger = log.New(os.Stdout, "[DEBUG] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
	l.infoLogger = log.New(os.Stdout, "[INFO] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
	l.warnLogger = log.New(os.Stdout, "[WARN] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
	l.errorLogger = log.New(os.Stdout, "[ERROR] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
	l.panicLogger = log.New(os.Stdout, "[PANIC] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
}

type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	panicLogger *log.Logger
	once        sync.Once
	level       LogLevel // 日志等级：colorlog.DebugLevel、colorlog.InfoLevel、colorlog.WarnLevel、colorlog.ErrorLevel
	colorful    bool     // 是否彩色
	col         *ColorType
}

func New() *Logger {
	return &Logger{level: InfoLevel}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) Color(col *ColorType) {
	l.col = col
}

func (l *Logger) Colorful(colorful bool) {
	l.colorful = colorful
}

func (l *Logger) logout(level LogLevel, format *string, args ...any) {
	l.once.Do(func() {
		l.init()
	})
	if l.level >= level {
		switch level {
		case DebugLevel:
			l.print(l.debugLogger, CyanBright, format, args...)
		case InfoLevel:
			l.print(l.infoLogger, GreenBright, format, args...)
		case WarnLevel:
			l.print(l.warnLogger, YellowBright, format, args...)
		case ErrorLevel:
			l.print(l.errorLogger, RedBright, format, args...)
		case PanicLevel:
			l.print(l.panicLogger, RedUnderLine, format, args...)
		}
	}
}

func (l *Logger) print(logger *log.Logger, col *ColorType, format *string, args ...any) {
	if l.colorful {
		if l.col != nil {
			// 设置过自定义颜色，优先自定义颜色
			if format != nil {
				_ = logger.Output(callDepth, string(*l.col)+fmt.Sprintf(*format, args...))
				return
			}
			_ = logger.Output(callDepth, string(*l.col)+fmt.Sprint(args...))
			return
		}
		// 没有自定义颜色，采用默认色系
		if format != nil {
			_ = logger.Output(callDepth, string(*col)+fmt.Sprintf(*format, args...))
			return
		}
		_ = logger.Output(callDepth, string(*col)+fmt.Sprint(args...))
		return
	}
	// 不输出颜色
	if format != nil {
		_ = logger.Output(callDepth, fmt.Sprintf(*format, args...))
		return
	}
	_ = logger.Output(callDepth, fmt.Sprintln(args...))
}

func (l *Logger) Debug(ctx context.Context, args ...any) {
	l.logout(DebugLevel, nil, args...)
}

func (l *Logger) Info(ctx context.Context, args ...any) {
	l.logout(InfoLevel, nil, args...)
}

func (l *Logger) Warn(ctx context.Context, args ...any) {
	l.logout(WarnLevel, nil, args...)
}

func (l *Logger) Error(ctx context.Context, args ...any) {
	l.logout(ErrorLevel, nil, args...)
}

func (l *Logger) Panic(ctx context.Context, args ...any) {
	l.logout(PanicLevel, nil, args...)
}

func (l *Logger) Debugf(ctx context.Context, fmt string, args ...any) {
	l.logout(DebugLevel, &fmt, args...)
}

func (l *Logger) Infof(ctx context.Context, fmt string, args ...any) {
	l.logout(InfoLevel, &fmt, args...)
}

func (l *Logger) Warnf(ctx context.Context, fmt string, args ...any) {
	l.logout(WarnLevel, &fmt, args...)
}

func (l *Logger) Errorf(ctx context.Context, fmt string, args ...any) {
	l.logout(ErrorLevel, &fmt, args...)
}

func (l *Logger) Panicf(ctx context.Context, fmt string, args ...any) {
	l.logout(PanicLevel, &fmt, args...)
}

type ColorType string

var (
	Reset = &reset
	// 标准
	White   = &white
	Red     = &red
	Green   = &green
	Yellow  = &yellow
	Blue    = &blue
	Magenta = &magenta
	Cyan    = &cyan
	// 高亮
	WhiteBright   = &whiteBright
	RedBright     = &redBright
	GreenBright   = &greenBright
	YellowBright  = &yellowBright
	BlueBright    = &blueBright
	MagentaBright = &magentaBright
	CyanBright    = &cyanBright
	// 斜体
	WhiteBevel   = &whiteBevel
	RedBevel     = &redBevel
	GreenBevel   = &greenBevel
	YellowBevel  = &yellowBevel
	BlueBevel    = &blueBevel
	MagentaBevel = &magentaBevel
	CyanBevel    = &cyanBevel
	// 下划线
	WhiteUnderLine   = &whiteUnderLine
	RedUnderLine     = &redUnderLine
	GreenUnderLine   = &greenUnderLine
	YellowUnderLine  = &yellowUnderLine
	BlueUnderLine    = &blueUnderLine
	MagentaUnderLine = &magentaUnderLine
	CyanUnderLine    = &cyanUnderLine
	// 背景色
	WhiteBg   = &whiteBg
	RedBg     = &redBg
	GreenBg   = &greenBg
	YellowBg  = &yellowBg
	BlueBg    = &blueBg
	MagentaBg = &magentaBg
	CyanBg    = &cyanBg
	// 删除线
	WhiteDelLine   = &whiteDelLine
	RedDelLine     = &redDelLine
	GreenDelLine   = &greenDelLine
	YellowDelLine  = &yellowDelLine
	BlueDelLine    = &blueDelLine
	MagentaDelLine = &magentaDelLine
	CyanDelLine    = &cyanDelLine

	// ================================================

	reset = ColorType([]byte{27, 91, 48, 109})
	// 标准
	white   = ColorType([]byte{27, 91, 51, 48, 109}) // 白色
	red     = ColorType([]byte{27, 91, 51, 49, 109}) // 红色
	green   = ColorType([]byte{27, 91, 51, 50, 109}) // 绿色
	yellow  = ColorType([]byte{27, 91, 51, 51, 109}) // 黄色
	blue    = ColorType([]byte{27, 91, 51, 52, 109}) // 蓝色
	magenta = ColorType([]byte{27, 91, 51, 53, 109}) // 紫色
	cyan    = ColorType([]byte{27, 91, 51, 54, 109}) // 青色
	// 高亮
	whiteBright   = ColorType([]byte{27, 91, 49, 59, 51, 48, 109})
	redBright     = ColorType([]byte{27, 91, 49, 59, 51, 49, 109})
	greenBright   = ColorType([]byte{27, 91, 49, 59, 51, 50, 109})
	yellowBright  = ColorType([]byte{27, 91, 49, 59, 51, 51, 109})
	blueBright    = ColorType([]byte{27, 91, 49, 59, 51, 52, 109})
	magentaBright = ColorType([]byte{27, 91, 49, 59, 51, 53, 109})
	cyanBright    = ColorType([]byte{27, 91, 49, 59, 51, 54, 109})
	// 斜体
	whiteBevel   = ColorType([]byte{27, 91, 51, 59, 51, 48, 109})
	redBevel     = ColorType([]byte{27, 91, 51, 59, 51, 49, 109})
	greenBevel   = ColorType([]byte{27, 91, 51, 59, 51, 50, 109})
	yellowBevel  = ColorType([]byte{27, 91, 51, 59, 51, 51, 109})
	blueBevel    = ColorType([]byte{27, 91, 51, 59, 51, 52, 109})
	magentaBevel = ColorType([]byte{27, 91, 51, 59, 51, 53, 109})
	cyanBevel    = ColorType([]byte{27, 91, 51, 59, 51, 54, 109})
	// 下划线
	whiteUnderLine   = ColorType([]byte{27, 91, 52, 59, 51, 48, 109})
	redUnderLine     = ColorType([]byte{27, 91, 52, 59, 51, 49, 109})
	greenUnderLine   = ColorType([]byte{27, 91, 52, 59, 51, 50, 109})
	yellowUnderLine  = ColorType([]byte{27, 91, 52, 59, 51, 51, 109})
	blueUnderLine    = ColorType([]byte{27, 91, 52, 59, 51, 52, 109})
	magentaUnderLine = ColorType([]byte{27, 91, 52, 59, 51, 53, 109})
	cyanUnderLine    = ColorType([]byte{27, 91, 52, 59, 51, 54, 109})
	// 背景色
	whiteBg   = ColorType([]byte{27, 91, 55, 59, 51, 48, 109})
	redBg     = ColorType([]byte{27, 91, 55, 59, 51, 49, 109})
	greenBg   = ColorType([]byte{27, 91, 55, 59, 51, 50, 109})
	yellowBg  = ColorType([]byte{27, 91, 55, 59, 51, 51, 109})
	blueBg    = ColorType([]byte{27, 91, 55, 59, 51, 52, 109})
	magentaBg = ColorType([]byte{27, 91, 55, 59, 51, 53, 109})
	cyanBg    = ColorType([]byte{27, 91, 55, 59, 51, 54, 109})
	// 删除线
	whiteDelLine   = ColorType([]byte{27, 91, 57, 59, 51, 48, 109})
	redDelLine     = ColorType([]byte{27, 91, 57, 59, 51, 49, 109})
	greenDelLine   = ColorType([]byte{27, 91, 57, 59, 51, 50, 109})
	yellowDelLine  = ColorType([]byte{27, 91, 57, 59, 51, 51, 109})
	blueDelLine    = ColorType([]byte{27, 91, 57, 59, 51, 52, 109})
	magentaDelLine = ColorType([]byte{27, 91, 57, 59, 51, 53, 109})
	cyanDelLine    = ColorType([]byte{27, 91, 57, 59, 51, 54, 109})
)

const (
	ErrorLevel LogLevel = iota + 1
	WarnLevel
	InfoLevel
	DebugLevel
	PanicLevel
)

const (
	callDepth = 4
)

type LogLevel int

var (
	logger = &Logger{level: InfoLevel}
)

func SetAppLogLevel(level LogLevel) {
	SetLevel(level)
}

func ConvertLevel(lvl string) LogLevel {
	lvl = strings.ToLower(lvl)
	switch lvl {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "panic":
		return PanicLevel
	default:
		return InfoLevel
	}
}

func AppLogLevel() LogLevel {
	return logger.level
}

func SetLevel(level LogLevel) {
	logger.level = level
}

func SetColorful(colorful bool) {
	logger.colorful = colorful
}

// 青色:Cyan
func Debug(ctx context.Context, args ...any) {
	logger.logout(DebugLevel, nil, args...)
}

// 青色:Cyan
func Debugf(ctx context.Context, format string, args ...any) {
	logger.logout(DebugLevel, &format, args...)
}

// 白色:White
func Info(ctx context.Context, args ...any) {
	logger.logout(InfoLevel, nil, args...)
}

// 白色:White
func Infof(ctx context.Context, format string, args ...any) {
	logger.logout(InfoLevel, &format, args...)
}

// 黄色:Yellow
func Warn(ctx context.Context, args ...any) {
	logger.logout(WarnLevel, nil, args...)
}

// 黄色:Yellow
func Warnf(ctx context.Context, format string, args ...any) {
	logger.logout(WarnLevel, &format, args...)
}

// 红色:Red
func Error(ctx context.Context, args ...any) {
	logger.logout(ErrorLevel, nil, args...)
}

// 红色:Red
func Errorf(ctx context.Context, format string, args ...any) {
	logger.logout(ErrorLevel, &format, args...)
}

func Panic(ctx context.Context, args ...any) {
	logger.logout(PanicLevel, nil, args...)
}

func Panicf(ctx context.Context, format string, args ...any) {
	logger.logout(PanicLevel, &format, args...)
}
