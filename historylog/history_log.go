//Package historylog for chatroom to log chat history
package historylog

import (
	"log"
	"os"
)

// HistoryLog log chat history
type HistoryLog struct {
	target string
	l      *log.Logger
	file   *os.File
}

// New return a history log
// need to target file
func New(logFile string) *HistoryLog {
	var customLog *log.Logger
	var file *os.File

	if logFile == "" {
		customLog = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		file, err := os.Create(logFile)
		if err != nil {
			log.Fatalln(err)
		}

		// this file need to be closed handy
		customLog = log.New(file, "", log.LstdFlags)
	}

	return &HistoryLog{
		target: logFile,
		l:      customLog,
		file:   file,
	}
}

// Close close the history log file
// if something went wrong it will return err
func (hl HistoryLog) Close() {
	log.Println("handy close history log")
	if hl.file != nil {
		if err := hl.file.Close(); err != nil {
			log.Println(err)
		}
	}
}

// PrintFrom print something about obj to log with newline
func (hl HistoryLog) PrintFrom(obj HistoryLogger) {
	hl.l.Println(obj.GenSummary())
}

// HistoryLogger can make a summary from itself
type HistoryLogger interface {
	GenSummary() string
}
