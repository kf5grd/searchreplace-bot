package logr

import "io"

type Logger struct {
	Writer      io.Writer // Where to write log messages
	EnableDebug bool      // Whether to write debug messages
	JSON        bool      // Whether to write messages in JSON format
}

func New(writer io.Writer, debug, json bool) *Logger {
	return &Logger{
		Writer:      writer,
		EnableDebug: debug,
		JSON:        json,
	}
}

type Level int

const (
	LevelUnknown Level = iota
	LevelDebug
	LevelInfo
	LevelError
)

var levelMap = map[Level]string{
	LevelUnknown: "UNKNOWN",
	LevelDebug:   "DEBUG",
	LevelInfo:    "INFO",
	LevelError:   "ERROR",
}

func (l Level) String() string {
	if s, ok := levelMap[l]; ok {
		return s
	}
	return levelMap[0]
}

type Msg struct {
	Time     int64  `json:"time"`
	FuncName string `json:"func_name"`
	Level    string `json:"level"`
	Message  string `json:"message"`
}
