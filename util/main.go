package util

import (
	"fmt"
	"log"
	"os"
)

// DebugLog is the file to write logs to
var DebugLog *os.File

// SetDualLogging means log statements will print to console and be saved in debug.log
func SetDualLogging() {
	var err error
	DebugLog, err = os.Create("debug.log")
	if err != nil {
		panic(err)
	}
}

// Errorf prints a formatted error to the ui console
func Errorf(fmtString string, fmtArgs ...interface{}) {
	log.Printf("[red]error:[lime] %v", fmt.Sprintf(fmtString, fmtArgs...))
	if DebugLog != nil {
		DebugLog.Sync()
	}
}

// Errorln prints text to the ui console
func Errorln(msg ...interface{}) {
	v := append([]interface{}{"[red]error:[lime]"}, msg...)
	log.Println(v...)
	if DebugLog != nil {
		DebugLog.Sync()
	}
}

// Warnf prints a formatted warning to the ui console
func Warnf(fmtString string, fmtArgs ...interface{}) {
	log.Printf("[orange]warn:[lime] %v", fmt.Sprintf(fmtString, fmtArgs...))
	if DebugLog != nil {
		DebugLog.Sync()
	}
}

// Warnln prints text to the ui console
func Warnln(msg ...interface{}) {
	v := append([]interface{}{"[orange]warn:[lime]"}, msg...)
	log.Println(v...)
	if DebugLog != nil {
		DebugLog.Sync()
	}
}

// Infof prints a formatted Infoing to the ui console
func Infof(fmtString string, fmtArgs ...interface{}) {
	log.Printf(fmtString, fmtArgs...)
	if DebugLog != nil {
		DebugLog.Sync()
	}
}

// Infoln prints text to the ui console
func Infoln(msg ...interface{}) {
	log.Println(msg...)
	if DebugLog != nil {
		DebugLog.Sync()
	}
}
