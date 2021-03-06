package system

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var infoLogger *log.Logger    // Info logger
var debugLogger *log.Logger   // Debug logger
var errorLogger *log.Logger   // Error logger
var warningLogger *log.Logger // Warning logger
var debugSwitch string        // Switch Debug on or off
var verboseSwitch string      // Switch verbose mode on or off

// calledFrom returns the name and the line number of the calling source file
func calledFrom() string {
	ret := ""
	_, file, no, ok := runtime.Caller(2)
	if ok {
		_, relfile := filepath.Split(file)
		ret = fmt.Sprintf("%s:%d: ", relfile, no)
	}
	return ret
}

// DebugLog sents text to the DebugLogWriter
func DebugLog(txt string, stuff ...interface{}) {
	if debugLogger != nil && debugSwitch == "1" {
		debugLogger.Printf(calledFrom()+txt+"\n", stuff...)
		fmt.Fprintf(os.Stderr, "DEBUG: "+txt+"\n", stuff...)
	}
}

// InfoLog sents text to the InfoLogWriter
func InfoLog(txt string, stuff ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(calledFrom()+txt+"\n", stuff...)
		if verboseSwitch == "on" {
			fmt.Fprintf(os.Stdout, "    INFO: "+txt+"\n", stuff...)
		}
	}
}

// WarningLog sents text to the WarningLogWriter
func WarningLog(txt string, stuff ...interface{}) {
	if warningLogger != nil {
		warningLogger.Printf(calledFrom()+txt+"\n", stuff...)
		if verboseSwitch == "on" {
			fmt.Fprintf(os.Stderr, "    WARNING: "+txt+"\n", stuff...)
		}
	}
}

// ErrorLog sents text to the ErrorLogWriter
func ErrorLog(txt string, stuff ...interface{}) error {
	if errorLogger != nil {
		errorLogger.Printf(calledFrom()+txt+"\n", stuff...)
		fmt.Fprintf(os.Stderr, "ERROR: "+txt+"\n", stuff...)
	}
	return fmt.Errorf(txt+"\n", stuff...)
}

// LogInit initialise the different log writer saptune will use
func LogInit(logFile, debug, verbose string) {
	var saptuneLog io.Writer
	//define log format
	logTimeFormat := time.Now().Format("2006-01-02 15:04:05.000 ")

	//create log file with desired read/write permissions
	//saptuneLog, err := os.OpenFile("/var/log/tuned/tuned.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	saptuneLog, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err.Error())
	}
	//saptuneWriter := io.MultiWriter(os.Stderr, saptuneLog)
	//log.SetOutput(saptuneWriter)

	//debugLogWriter := io.MultiWriter(os.Stderr, saptuneLog)
	//infoLogWriter := io.MultiWriter(os.Stdout, saptuneLog)
	//warningLogWriter := io.MultiWriter(os.Stderr, saptuneLog)
	//errorLogWriter := io.MultiWriter(os.Stderr, saptuneLog)

	//debugLogger = log.New(debugLogWriter, logTimeFormat+"DEBUG    saptune.", 0)
	//infoLogger = log.New(infoLogWriter, logTimeFormat+"INFO     saptune.", 0)
	//warningLogger = log.New(warningLogWriter, logTimeFormat+"WARNING  saptune.", 0)
	//errorLogger = log.New(errorLogWriter, logTimeFormat+"ERROR    saptune.", 0)
	//errorLogger = log.New(errorLogWriter, logTimeFormat+"ERROR    saptune.", log.Lshortfile)
	//log.SetFlags(0)

	debugLogger = log.New(saptuneLog, logTimeFormat+"DEBUG    saptune.", 0)
	infoLogger = log.New(saptuneLog, logTimeFormat+"INFO     saptune.", 0)
	warningLogger = log.New(saptuneLog, logTimeFormat+"WARNING  saptune.", 0)
	errorLogger = log.New(saptuneLog, logTimeFormat+"ERROR    saptune.", 0)

	debugSwitch = debug
	verboseSwitch = verbose
}
