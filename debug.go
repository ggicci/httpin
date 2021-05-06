package httpin

import (
	"log"
)

var (
	debugOn = false
)

func DebugOn() {
	debugOn = true
}

func DebugOff() {
	debugOn = false
}

func debug(format string, v ...interface{}) {
	if debugOn {
		log.Printf(format, v...)
	}
}
