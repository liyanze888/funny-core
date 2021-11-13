package fn_log

import (
	"fmt"
	"github.com/liyanze888/funny-core/fn_utils"
	"log"
	"os"
)

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
const Deep = 2

func Print(v ...interface{}) {
	Output(Deep, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Output(Deep, fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	log.Output(Deep, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	Output(Deep, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	Output(Deep, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	Output(Deep, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	Output(Deep, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	Output(Deep, s)
	panic(s)
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is the count of the number of
// frames to skip when computing the file name and line number
// if Llongfile or Lshortfile is set; a value of 1 will print the details
// for the caller of Output.
func Output(calldepth int, s string) error {
	callId, b := fn_utils.Get("CallId")
	StreamCallId, b2 := fn_utils.Get("StreamCallId")
	if b && b2 {
		s = fmt.Sprintf("[CallId = %v-%v] %s", callId, StreamCallId, s)
	} else if b {
		s = fmt.Sprintf("[CallId = %v] %s", callId, s)
	} else if b2 {
		s = fmt.Sprintf("[CallId = %v] %s", StreamCallId, s)
	}

	return log.Output(calldepth+1, s) // +1 for this frame.
}
