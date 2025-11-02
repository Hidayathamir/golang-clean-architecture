package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/log"
)

func EmitLogEntry(entry *logrus.Entry) {
	ctx := entry.Context
	if ctx == nil {
		ctx = context.Background()
	}

	if logger == nil {
		return
	}

	severity, severityText := convertSeverity(entry.Level)

	if !logger.Enabled(ctx, log.EnabledParameters{Severity: severity}) {
		return
	}

	var record log.Record

	if !entry.Time.IsZero() {
		record.SetTimestamp(entry.Time)
	} else {
		record.SetTimestamp(time.Now())
	}

	record.SetSeverity(severity)
	record.SetSeverityText(severityText)
	record.SetBody(log.StringValue(entry.Message))

	for key, value := range entry.Data {
		record.AddAttributes(log.String(key, JSONStr(value)))
	}

	logger.Emit(ctx, record)
}

var limitChar = 10000

func JSONStr(v any) string {
	jsonByte, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	jsonStr := string(jsonByte)
	if len(jsonStr) > limitChar {
		jsonStr = jsonStr[:limitChar] + "..."
		return jsonStr
	}
	return jsonStr
}

func convertSeverity(level logrus.Level) (log.Severity, string) {
	switch level {
	case logrus.TraceLevel:
		return log.SeverityTrace, "TRACE"
	case logrus.DebugLevel:
		return log.SeverityDebug, "DEBUG"
	case logrus.InfoLevel:
		return log.SeverityInfo, "INFO"
	case logrus.WarnLevel:
		return log.SeverityWarn, "WARN"
	case logrus.ErrorLevel:
		return log.SeverityError, "ERROR"
	case logrus.FatalLevel:
		return log.SeverityFatal3, "FATAL"
	case logrus.PanicLevel:
		return log.SeverityFatal4, "PANIC"
	default:
		return log.SeverityInfo, strings.ToUpper(level.String())
	}
}
