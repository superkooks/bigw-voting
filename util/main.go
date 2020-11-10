package util

import (
	"fmt"
	"log"
)

// Errorf prints a formatted error to the ui console
func Errorf(fmtString string, fmtArgs ...interface{}) {
	log.Printf("[red]error:[lime] %v", fmt.Sprintf(fmtString, fmtArgs...))
}

// Errorln prints text to the ui console
func Errorln(msg ...interface{}) {
	v := append([]interface{}{"[red]error:[lime]"}, msg...)
	log.Println(v...)
}

// Warnf prints a formatted warning to the ui console
func Warnf(fmtString string, fmtArgs ...interface{}) {
	log.Printf("[orange]warn:[lime] %v", fmt.Sprintf(fmtString, fmtArgs...))
}

// Warnln prints text to the ui console
func Warnln(msg ...interface{}) {
	v := append([]interface{}{"[orange]warn:[lime]"}, msg...)
	log.Println(v...)
}

// Infof prints a formatted Infoing to the ui console
func Infof(fmtString string, fmtArgs ...interface{}) {
	log.Printf(fmtString, fmtArgs...)
}

// Infoln prints text to the ui console
func Infoln(msg ...interface{}) {
	log.Println(msg...)
}

// // Printf prints text without the time
// func Printf(fmtString string, fmtArgs ...interface{}) {
// 	ui.GetConsoleWriter().Write([]byte(fmt.Sprintf(fmtString, fmtArgs...)))
// }

// // Println prints text without the time
// func Println(msg ...interface{}) {
// 	ui.GetConsoleWriter().Write([]byte(fmt.Sprintln(msg...)))
// }
