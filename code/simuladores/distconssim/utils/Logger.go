package utils

import (
	"log"
	"os"
)

type Logger struct {
	// Logs
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func InitLoggers(exp string, name string) *Logger {

	// Initialize log
	fLog, err := os.OpenFile(LoggerPath+exp+"/Log_"+name+".log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	myLogger := Logger{}
	myLogger.Trace = log.New(fLog,
		"TRACE: \t\t["+name+"] ", log.Ltime|log.Lmicroseconds|log.Lshortfile)

	myLogger.Info = log.New(fLog,
		"INFO: \t\t["+name+"] ", log.Ltime|log.Lmicroseconds|log.Lshortfile)

	myLogger.Warning = log.New(fLog,
		"WARNING: \t["+name+"] ", log.Ltime|log.Lmicroseconds|log.Lshortfile)

	myLogger.Error = log.New(fLog,
		"ERROR: \t\t["+name+"] ", log.Ltime|log.Lmicroseconds|log.Lshortfile)

	return &myLogger
}
