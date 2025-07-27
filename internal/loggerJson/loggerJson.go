package loggerjson

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type Logger struct {
	output   io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {

	return &Logger{
		minLevel: minLevel,
		output:   out,
	}
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

func (log *Logger) PrintDebug(message string, properties map[string]string) {
	log.Print(LevelDebug, message, properties)
}

func (log *Logger) PrintInfo(message string, properties map[string]string) {
	log.Print(LevelInfo, message, properties)
}

func (log *Logger) PrintWarn(message string, properties map[string]string) {
	log.Print(LevelWarn, message, properties)

}
func (log *Logger) PrintError(err error, properties map[string]string) {
	log.Print(LevelWarn, err.Error(), properties)
}
func (log *Logger) PrintFatal(err error, properties map[string]string) {
	log.Print(LevelWarn, err.Error(), properties)
	os.Exit(1)
}

func (log *Logger) Write(message []byte) (n int, err error) {
	return log.Print(LevelError, string(message), nil)
}

func (log *Logger) Print(level Level, message string, properties map[string]string) (int, error) {

	if level < log.minLevel {
		return 0, nil
	}

	auxlog := struct {
		Level      Level             `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level,
		Time:       time.Now().String(),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		auxlog.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(auxlog)

	if err != nil {
		line = []byte(LevelError.String() + ":unable to marshall errror json" + err.Error())
	}

	log.mu.Lock()
	defer log.mu.Unlock()

	return log.output.Write(append(line, '\n'))
}
