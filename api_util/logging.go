// This package provides various utility functions for working with interfaces, maps, 
// arrays, and JSON/YAML conversion. Each function has a comment explaining its 
// purpose and usage. Let me know if you need further clarification on any specific function!

package api_util


import (
	"fmt"
	"io"
	"log"
	"os"
)

// Test results constants
const (
	Passed         = "Passed"
	Failed         = "Failed"
	Skipped        = "Skipped"
	SchemaMismatch = "SchemaMismatch"
	Total          = "Total"
)

// Colors for better logging
const (
	RED    = "\033[1;31m"
	GREEN  = "\033[1;32m"
	YELLOW = "\033[1;33m"
	BLUE   = "\033[1;34m"
	AQUA   = "\033[1;36m"
	END    = "\033[0m"
)

// NewLogger creates a new logger with the specified output writer.
// The logger is configured with the current date, time in microseconds,
// and the file name and line number where the log statement is called.
// The output writer is used to write the log messages.
// 
// Parameters:
//   - out: The output writer to write the log messages to.
// 
// Returns:
//   - *log.Logger: The created logger.
func NewLogger(out io.Writer) *log.Logger {
	Logger = log.New(out, "", (log.Ldate | log.Lmicroseconds | log.Lshortfile))
	return Logger
}

// NewStdLogger returns a new instance of log.Logger that writes log messages to os.Stdout.
func NewStdLogger() *log.Logger {
	return NewLogger(os.Stdout)
}

// NewFileLogger creates a new file logger with the specified file path.
// It opens the file at the given path with read, write, and create permissions.
// If the file cannot be opened, it returns nil and prints an error message.
// The returned logger can be used to write log messages to the file.
func NewFileLogger(path string) *log.Logger {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Can't open %s, err: %s", path, err.Error())
		return nil
	}
	return NewLogger(f)
}

// There is only one logger per process.
var Logger *log.Logger

// Whether verbose mose is on
var Verbose bool
