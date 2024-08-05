package logs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type StructuredLogger struct {
	Logger *logrus.Logger
}

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

func (s *StructuredLogger) NewLogEntry(req *http.Request) middleware.LogEntry {
	logEntry := &StructuredLoggerEntry{Logger: logrus.NewEntry(s.Logger)}
	loggerFields := logrus.Fields{}

	if reqId := middleware.GetReqID(req.Context()); reqId != "" {
		loggerFields["request_id"] = reqId
	}

	loggerFields["http_method"] = req.Method
	loggerFields["remote_address"] = req.RemoteAddr
	loggerFields["uri"] = req.RequestURI

	logEntry.Logger = logEntry.Logger.WithFields(loggerFields)
	logEntry.Logger.Info("Request started")

	return logEntry
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

func LogEntry(req *http.Request) logrus.FieldLogger {
	logEntry := middleware.GetLogEntry(req)

	return logEntry.(*StructuredLoggerEntry).Logger
}

func (s *StructuredLoggerEntry) Write(respStatus, respBytes int, httpHeader http.Header, elapsedTime time.Duration, extraData any) {
	s.Logger = s.Logger.WithFields(logrus.Fields{
		"response_bytes_length": respBytes,
		"response_elapsed_time": elapsedTime.Round(time.Millisecond / 100).String(),
		"response_status":       respStatus,
	})

	s.Logger.Info("Request completed")
}

func (s *StructuredLoggerEntry) Panic(v any, stack []byte) {
	s.Logger = s.Logger.WithFields(logrus.Fields{
		"panic": fmt.Sprintf("%+v", v),
		"stack": string(stack),
	})
}
