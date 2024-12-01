package logger

import (
	"log"
)

var enabled bool

func init() {
	enabled = true
}

func Enable() {
	enabled = true
}

func Disable() {
	enabled = false
}

func Info(format string, args ...interface{}) {
	if enabled {
		log.Printf("[INFO] "+format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	if enabled {
		log.Printf("[WARN] "+format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if enabled {
		log.Printf("[ERROR] "+format, args...)
	}
}
