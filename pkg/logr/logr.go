package logr

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// getFrame and getCaller taken from https://stackoverflow.com/questions/35212985/is-it-possible-get-information-about-caller-function-in-golang
func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

func getCaller() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	frame := strings.Split(getFrame(2).Function, ".")
	return frame[len(frame)-1]
}

func (l *Logger) Write(name string, level Level, s string, a ...interface{}) {
	var (
		now        = time.Now().UTC()
		timeFormat = "02Jan2006 15:04:05"
	)

	if l.JSON {
		msg := Msg{
			Time:     now.Unix(),
			FuncName: name,
			Level:    level.String(),
			Message:  fmt.Sprintf(s, a...),
		}
		fmt.Fprintf(l.Writer, toJson(msg)+"\n")
		return
	}

	a = append([]interface{}{strings.ToUpper(now.Format(timeFormat)), name, level}, a...)
	fmt.Fprintf(l.Writer, "[%v][%s] %v: "+s+"\n", a...)
}

func (l *Logger) Error(s string, a ...interface{}) {
	l.Write(getCaller(), LevelError, s, a...)
}

func (l *Logger) Info(s string, a ...interface{}) {
	l.Write(getCaller(), LevelInfo, s, a...)
}

func (l *Logger) Debug(s string, a ...interface{}) {
	if !l.EnableDebug {
		return
	}
	l.Write(getCaller(), LevelDebug, s, a...)
}
