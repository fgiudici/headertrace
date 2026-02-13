package logging

import (
	"log"
	"os"
	"strings"
)

type level int

const (
	TRACE level = iota
	DEBUG
	INFO
	WARN
	ERROR
)

var lvl = INFO

// Init configures the logger. It reads LOG_LEVEL from the environment
// (one of: TRACE, DEBUG, INFO, WARN, ERROR) and sets a simple prefix.
func Init() {
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		switch strings.ToUpper(v) {
		case "TRACE":
			lvl = TRACE
		case "DEBUG":
			lvl = DEBUG
		case "INFO":
			lvl = INFO
		case "WARN":
			lvl = WARN
		case "ERROR":
			lvl = ERROR
		}
	}
	log.SetFlags(log.LstdFlags)
}

func Tracef(format string, v ...interface{}) {
	if lvl <= TRACE {
		log.Printf("TRACE: "+format, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if lvl <= DEBUG {
		log.Printf("DEBUG: "+format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	if lvl <= INFO {
		log.Printf("INFO: "+format, v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if lvl <= WARN {
		log.Printf("WARN: "+format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if lvl <= ERROR {
		log.Printf("ERROR: "+format, v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	log.Printf("FATAL: "+format, v...)
	os.Exit(1)
}
