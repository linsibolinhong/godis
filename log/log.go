package log

import "fmt"

type LogLevel int;

const (
	LOG_TRACE = LogLevel(iota)
	LOG_DEBUG
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_CRIT
	LOG_OFF
)

var levelSlice = []string{}

func init() {
	levelSlice = make([]string, LOG_OFF + 1)

	levelSlice[LOG_TRACE] = "trace"
	levelSlice[LOG_DEBUG] = "debug"
	levelSlice[LOG_INFO] = "info"
	levelSlice[LOG_WARN] = "warn"
	levelSlice[LOG_ERROR] = "error"
	levelSlice[LOG_CRIT] = "critical"
}

func logPrint(level LogLevel, f string, a ...interface{}) {
	fmt.Printf("[%v]" + f + "\n", levelSlice[level], a)
}

func Trace(f string, a ...interface{}) {
	logPrint(LOG_TRACE, f, a)
}

func Debug(f string, a ...interface{}) {
	logPrint(LOG_DEBUG, f, a)
}

func Info(f string, a ...interface{}) {
	logPrint(LOG_INFO, f, a)
}

func Warn(f string, a ...interface{}) {
	logPrint(LOG_WARN, f, a)
}

func Error(f string, a ...interface{}) {
	logPrint(LOG_ERROR, f, a)
}

func Critical(f string, a ...interface{}) {
	logPrint(LOG_CRIT, f, a)
}
