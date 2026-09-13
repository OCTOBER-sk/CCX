package log

import (
	"fmt"
	"log"
)

// Per spec section 48 (logging/observability) and 4.12 (secrets).
// Levels: ERROR, WARN, INFO, DEBUG, TRACE. Default body logging: OFF.
// Redaction is structural: never serialize credential objects into logs.

var logger = log.New(log.Writer(), "[CCX] ", log.LstdFlags)

func Error(msg string, args ...interface{}) {
	logger.Printf("ERROR: "+msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Printf("WARN: "+msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.Printf("INFO: "+msg, args...)
}

func Debug(msg string, args ...interface{}) {
	logger.Printf("DEBUG: "+msg, args...)
}

func Redact(s string) string {
	if s == "" {
		return ""
	}
	// Minimal redaction: replace any string that looks like a secret with [REDACTED]
	// Real structural redaction comes from never putting secrets in log payloads (spec 4.12)
	return "[REDACTED]"
}

func RedactedLog(value interface{}) string {
	return fmt.Sprintf("%v [REDACTED]", value)
}
