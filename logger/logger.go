package logger

import (
	"log"
	"os"
)

var Error *log.Logger = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var Info *log.Logger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

func Infof(format string, v ...any) {
	Info.Printf(format, v...)
}

func Infoln(v ...any) {
	Info.Println(v...)
}

func Errorf(format string, v ...any) {
	Error.Printf(format, v...)
}

func Errorln(v ...any) {
	Error.Println(v...)
}
