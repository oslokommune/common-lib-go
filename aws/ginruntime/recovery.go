package ginruntime

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type GoRoutine = string

type Frame struct {
	File     string
	Line     string
	Function string
}

type StackTrace struct {
	GoRoutine string
	Stack     []Frame
	Reason    any
}

// Extracts stack trace from a `StackTrace` error and returns it
// Usage:
// ```go
// zerolog.ErrorStackMarshaler = StackTraceMarshaller
// goroutine, stack := GetStack()
// stacktrace := StackTrace{GoRoutine: goroutine, Stack: stack, Reason: r}
// log.Info().Err(stacktrace).Msg("An error occurred")
// ```
func StackTraceMarshaller(err error) any {
	panicErr, ok := err.(StackTrace)
	if !ok {
		log.Warn().Err(err).Msg("Cannot extract stack trace from non-StackTrace error")
		return nil
	}

	out := make([]map[string]string, len(panicErr.Stack))
	for i, frame := range panicErr.Stack {
		out[i] = map[string]string{
			"file":     frame.File,
			"line":     frame.Line,
			"function": frame.Function,
		}
	}
	return out
}

// Gin middleware for recovering from panics and logging to zerolog.
// Assumes `StackTraceMarshaller` is being used.
//
//	Usage:
//	```go
//		engine.Use(RecoveryMiddleware)
//	```
func RecoveryMiddleware(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			goroutine, stack := GetStack()
			stacktrace := StackTrace{GoRoutine: goroutine, Stack: stack, Reason: r}
			stacktrace = stacktrace.SkipFramesAfterPanic()
			log.Error().Stack().Err(stacktrace).Msg("A panic occurred, which will cause a 500 INTERNAL_SERVER_ERROR response")
			c.AbortWithStatus(500)
		}
	}()
	c.Next()
}

func (s StackTrace) Error() string {
	return fmt.Sprintf("panic in %s: %v+", s.GoRoutine, s.Reason)
}

func (s StackTrace) Skip(n int) StackTrace {
	if n < 0 {
		n = 0
	}
	if n >= len(s.Stack) {
		n = len(s.Stack) - 1
	}

	return StackTrace{
		GoRoutine: s.GoRoutine,
		Stack:     s.Stack[n:],
		Reason:    s.Reason,
	}
}

func (s StackTrace) SkipFramesAfterPanic() StackTrace {
	return s.Skip(s.estimateNumberOfInternalFrames())
}

func (s StackTrace) estimateNumberOfInternalFrames() int {
	panicFilePattern := regexp.MustCompile(`^.*runtime/panic.go$`)
	panicFuncPattern := regexp.MustCompile(`^panic\((.*)\)$`)
	panicFrameIndex := 0
	for i := len(s.Stack) - 1; i >= 0; i-- {
		frame := s.Stack[i]
		if panicFuncPattern.MatchString(frame.Function) && panicFilePattern.MatchString(frame.File) {
			panicFrameIndex = i
			break
		}
	}

	if panicFrameIndex+1 < len(s.Stack) {
		return panicFrameIndex + 1
	}
	return 0
}

// Extracts the current goroutine and stack.
// The stack includes the `GetStack` function call and its call to `debug.Stack()`.
func GetStack() (GoRoutine, []Frame) {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	goroutine, lines := lines[0], lines[1:]
	goroutine = strings.TrimSuffix(goroutine, ":")
	frames := make([]Frame, len(lines)/2)

	for i := 0; i < len(frames); i++ {

		function := lines[2*i]
		location := lines[2*i+1]
		location = strings.TrimSpace(location)
		locationParts := strings.Split(location, ":")

		line := ""
		file := locationParts[0]
		if 1 < len(locationParts) {
			// Strip away the frame offset ("+0x2a" etc.)
			line = strings.Split(locationParts[1], " ")[0]
		}
		frames[i] = Frame{Function: function, File: file, Line: line}
	}
	return goroutine, frames
}
