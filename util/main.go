package util

import (
	"fmt"
	"log"
	"os"
)

var debugLog *os.File

// SetDualLogging means log statements will print to console and be saved in debug.log
func SetDualLogging() {
	var err error
	debugLog, err = os.Create("debug.log")
	if err != nil {
		panic(err)
	}
}

// Errorf prints a formatted error to the ui console
func Errorf(fmtString string, fmtArgs ...interface{}) {
	log.Printf("[red]error:[lime] %v", fmt.Sprintf(fmtString, fmtArgs...))
	if debugLog != nil {
		debugLog.WriteString(fmt.Sprint("error: ", fmt.Sprintf(fmtString, fmtArgs...)))
		debugLog.Sync()
	}
}

// Errorln prints text to the ui console
func Errorln(msg ...interface{}) {
	v := append([]interface{}{"[red]error:[lime]"}, msg...)
	log.Println(v...)
	if debugLog != nil {
		debugLog.WriteString(fmt.Sprintln(append([]interface{}{"error:"}, msg...)...))
		debugLog.Sync()
	}
}

// Warnf prints a formatted warning to the ui console
func Warnf(fmtString string, fmtArgs ...interface{}) {
	log.Printf("[orange]warn:[lime] %v", fmt.Sprintf(fmtString, fmtArgs...))
	if debugLog != nil {
		debugLog.WriteString(fmt.Sprint("warn: ", fmt.Sprintf(fmtString, fmtArgs...)))
		debugLog.Sync()
	}
}

// Warnln prints text to the ui console
func Warnln(msg ...interface{}) {
	v := append([]interface{}{"[orange]warn:[lime]"}, msg...)
	log.Println(v...)
	if debugLog != nil {
		debugLog.WriteString(fmt.Sprintln(append([]interface{}{"warn:"}, msg...)...))
		debugLog.Sync()
	}
}

// Infof prints a formatted Infoing to the ui console
func Infof(fmtString string, fmtArgs ...interface{}) {
	log.Printf(fmtString, fmtArgs...)
	if debugLog != nil {
		debugLog.WriteString(fmt.Sprintf(fmtString, fmtArgs...))
		debugLog.Sync()
	}
}

// Infoln prints text to the ui console
func Infoln(msg ...interface{}) {
	log.Println(msg...)
	if debugLog != nil {
		_, err := debugLog.WriteString(fmt.Sprintln(msg...))
		if err != nil {
			panic(err)
		}
		debugLog.Sync()
	}
}
