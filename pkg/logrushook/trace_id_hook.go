package logrushook

import (
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/traceidctx"

	"github.com/sirupsen/logrus"
)

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
	traceID := traceidctx.Get(ctx)
	if traceID != "" {
		entry.Data["trace_id"] = traceID
	}
	return nil
}
