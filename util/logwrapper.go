package util

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
	var baseLogger = logrus.New()

	var standardLogger = &StandardLogger{baseLogger}

	standardLogger.Formatter = &logrus.JSONFormatter{}

	return standardLogger
}

// Declare variables to store log messages as new Events
var (
	badRequestError = Event{400, "Bad Request: %s"}
	internalServerError = Event{500, "Internal Server Error: %s"}
)

// Define standard error messages

func (l *StandardLogger) InternalServerError(argumentName string) {
	l.Errorf(internalServerError.message, argumentName)
}

func (l *StandardLogger) BadRequestError(argumentName string) {
	l.Errorf(badRequestError.message, argumentName)
}

func LogRequest(l *StandardLogger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}