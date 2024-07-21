package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"os"

	"github.com/fatih/color"
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var levelStrings = []string{
	"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL",
}

func colorForLevel(level int) color.Attribute {
	switch level {
	case TRACE:
		return color.FgCyan
	case DEBUG:
		return color.FgBlue
	case INFO:
		return color.FgGreen
	case WARN:
		return color.FgYellow
	case ERROR:
		return color.FgRed
	case FATAL:
		return color.FgHiRed
	default:
		return color.Reset
	}
}

type logEvent struct {
	time     time.Time
	file     string
	line     int
	level    int
	message  string
	function string
}

type eventInfo struct {
	*logEvent
}

func (e *eventInfo) Time() time.Time {
	return e.time
}

func (e *eventInfo) File() string {
	return e.file
}

func (e *eventInfo) Line() int {
	return e.line
}

func (e *eventInfo) Level() int {
	return e.level
}

func (e *eventInfo) Message() string {
	return e.message
}

func (e *eventInfo) Function() string {
	return e.function
}

type EventInfo interface {
	Time() time.Time
	File() string
	Line() int
	Level() int
	Message() string
	Function() string
}
type logFunc func(EventInfo)
type Logger struct {
	level     int
	quiet     bool
	useColor  bool
	callbacks []logFunc
	mu        sync.Mutex
}

func NewLogger(level int, quiet bool) *Logger {
	return &Logger{level, quiet, true, make([]logFunc, 0), sync.Mutex{}}
}

func (l *Logger) AddCallback(callback logFunc) {
	l.mu.Lock()
	l.callbacks = append(l.callbacks, callback)
	l.mu.Unlock()
}

func (l *Logger) SetUseColor(useColor bool) {
	l.mu.Lock()
	l.useColor = useColor
	l.mu.Unlock()
}
func (l *Logger) log(level int, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if level < l.level || l.quiet {
		return
	}

	now := time.Now()
	pc, file, line, _ := runtime.Caller(2)
	function := runtime.FuncForPC(pc).Name()

	event := logEvent{
		time:     now,
		file:     filepath.Base(file),
		line:     line,
		level:    level,
		message:  fmt.Sprintf(format, args...),
		function: function,
	}

	if len(l.callbacks) > 0 {
		l.handleWithCallbacks(&event)
	} else {
		l.printEvent(&event)
	}
}

func (l *Logger) handleWithCallbacks(event *logEvent) {
	info := &eventInfo{event}
	for _, callback := range l.callbacks {
		if l.useColor {
			color.Set(colorForLevel(event.level))
			defer color.Unset()
		}
		callback(info)
	}
}
func (l *Logger) printEvent(event *logEvent) {
	if l.useColor {
		color.Set(colorForLevel(event.level))
		defer color.Unset()
	}
	timestamp := event.time.Format("2006-01-02 15:04:05")
	levelStr := levelStrings[event.level]

	fmt.Printf("%s %-5s %s:%d %s: %s\n", timestamp, levelStr, event.file, event.line, event.function, event.message)
}

func (l *Logger) Trace(format string, args ...interface{}) {
	l.log(TRACE, format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}
