package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func Println(v ...interface{}) {
	// Capture the call stack
	callers := make([]uintptr, 10)
	n := runtime.Callers(2, callers[:])
	frames := runtime.CallersFrames(callers[:n])

	// Iterate through the call stack
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		// Read the source code of the caller
		sourceLine := getSourceLine(frame.File, frame.Line)

		// Log the details
		fmt.Fprintf(os.Stderr, "%s:%s:%d: %s\n", time.Now().Format(time.RFC3339), frame.File, frame.Line, frame.Function)
		fmt.Fprintf(os.Stderr, "\n%s\n", sourceLine)
	}

	// Log the provided message
	for _, vv := range v {
		fmt.Fprintf(os.Stderr, "%s: %v\n", time.Now().Format(time.RFC3339), vv)
	}
}

func getSourceLine(filename string, line int) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sourceLines []string
	currentLine := 1
	for scanner.Scan() {
		if currentLine >= line-5 && currentLine <= line+5 {
			lineText := scanner.Text()
			if currentLine == line {
				lineText = fmt.Sprintf("\033[31m%d: %s\033[0m", currentLine, lineText) // Red color
			} else {
				lineText = fmt.Sprintf("%d: %s", currentLine, lineText)
			}
			sourceLines = append(sourceLines, lineText)
		}
		currentLine++
	}
	return strings.Join(sourceLines, "\n")
}
