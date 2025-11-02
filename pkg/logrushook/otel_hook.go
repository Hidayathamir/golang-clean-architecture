package logrushook

import (
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ logrus.Hook = &OtelHook{}

type OtelHook struct{}

func NewOtelHook() *OtelHook {
	return &OtelHook{}
}

func (h *OtelHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *OtelHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	traceID := telemetry.GetTraceID(ctx)
	if traceID != "" {
		entry.Data["trace_id"] = traceID
	}

	spanID := telemetry.GetSpanID(ctx)
	if spanID != "" {
		entry.Data["span_id"] = spanID
	}

	return nil
}
