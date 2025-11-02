package logrushook

import (
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ logrus.Hook = &TraceID{}

type TraceID struct{}

func NewTraceID() *TraceID {
	return &TraceID{}
}

func (h *TraceID) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *TraceID) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}
	traceID := telemetry.GetTraceID(ctx)
	if traceID != "" {
		entry.Data["trace_id"] = traceID
	}
	return nil
}
