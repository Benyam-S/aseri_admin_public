package log

import (
	"fmt"
	"os"
	"time"

	"github.com/Benyam-S/aseri/entity"
)

// Logger is a type that defines a logger type
type Logger struct {
	ServerLogFile string
	BotLogFile    string
}

// Log is a method that will log the given statement to the selected log file
func (l *Logger) Log(stmt, stmtType, logFile string) {

	if logFile == entity.BotLogFile {
		logFile = l.BotLogFile
	} else {
		logFile = l.ServerLogFile
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()

		if stmtType == "w" {
			stmtType = "Warning:"

		} else if stmtType == "e" {
			stmtType = "Error:"

		}

		stmt = fmt.Sprintf("%s %s  At %s", stmtType, stmt, time.Now())

		fmt.Fprintln(file, stmt)

	}
}

// LogToParent is a method that will log the given statement to the program starter
func (l *Logger) LogToParent(stmt, stmtType string) {

	if stmtType == "w" {
		stmtType = "Warning:"

	} else if stmtType == "e" {
		stmtType = "Error:"

	} else {
		stmtType = "Success:"

	}

	stmt = fmt.Sprintf("%s %s", stmtType, stmt)
	fmt.Println(stmt)
}

// LogFileError is a method that will log the given statement as an error to the selected log file
func (l *Logger) LogFileError(stmt, logFile string) {

	if logFile == entity.BotLogFile {
		logFile = l.BotLogFile
	} else {
		logFile = l.ServerLogFile
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()

		stmt = fmt.Sprintf("Error: %s  At %s", stmt, time.Now())

		fmt.Fprintln(file, stmt)

	}
}
