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
func New(logFile string) HistoryLog {
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

	return HistoryLog{
		target: logFile,
		l:      customLog,
		file:   file,
	}
}

// PrintFrom print something about obj to log with newline
func (hl HistoryLog) PrintFrom(msg string, obj HistoryLogger) {
	hl.l.Println(msg, "from", obj.GenSummary())
}

// LoginFrom print something about obj to log with newline and exit
func (hl HistoryLog) LoginFrom(obj HistoryLogger) {
	hl.l.Println("[LOGIN]", obj.GenSummary())
}

// QuitFrom print something about obj to log with newline and exit
func (hl HistoryLog) QuitFrom(obj HistoryLogger) {
	hl.l.Println("[QUIT]", obj.GenSummary())
}

// HistoryLogger can make a summary from itself
type HistoryLogger interface {
	GenSummary() string
}
