package telemetry

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/caller"
	"github.com/gofiber/fiber/v2"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/trace"
)

func StartController(ctx *fiber.Ctx) trace.Span {
	newctx, span := tracer.Start(ctx.UserContext(), caller.FuncName(caller.WithSkip(1)))
	ctx.SetUserContext(newctx)
	return span
}

func StartConsumer(ctx context.Context, record *kgo.Record, skips ...int) (context.Context, trace.Span) {
	// record.Context contains the span context extracted from Kafka headers by kotel.
	// We want to start a child span, but use the provided 'ctx' (shutdown context)
	// as the base context so that cancellation is respected.
	remoteSpanCtx := trace.SpanContextFromContext(record.Context)
	if remoteSpanCtx.IsValid() {
		ctx = trace.ContextWithRemoteSpanContext(ctx, remoteSpanCtx)
	}

	skip := 1
	if len(skips) > 0 {
		skip += skips[0]
	}

	return tracer.Start(ctx, caller.FuncName(caller.WithSkip(skip)))
}

// StartConsumerBatch implements a "Batch + Bridge" tracing pattern for Kafka consumers.
//
// Pattern Architecture:
//  1. Batch Span: Represents the overall execution of the consumer loop. It links to
//     every message's producer context to show what the batch "contains".
//  2. Bridge Spans (Inner): For each message, it starts a span that is a direct child
//     of the message's original producer. This preserves the individual message lineage.
//  3. Correlation Link: Each Bridge Span is linked back to the Batch Span, allowing
//     navigating from a single message's trace back to the batch that processed it.
//
// This approach is preferred over a strict parent-child hierarchy for large batches
// to keep trace visualizations readable while maintaining full bi-directional correlation.
//
// That is ai documentation, in my language, we connecting via trace.Link from
// original trace (in producer) connected to new trace (in consumer that read batch).
func StartConsumerBatch(originalCtx context.Context, records []*kgo.Record) (context.Context, trace.Span) {
	// Step 1: Collect all remote span contexts from the batch records.
	var links []trace.Link
	for _, record := range records {
		remoteSpanCtx := trace.SpanContextFromContext(record.Context)
		if remoteSpanCtx.IsValid() {
			links = append(links, trace.Link{SpanContext: remoteSpanCtx})
		}
	}

	// Create the Batch Span as an anchor for the whole operation, linking all producers.
	ctx, span := tracer.Start(originalCtx, caller.FuncName(caller.WithSkip(1)), trace.WithLinks(links...))

	// Step 2: Create "Bridge Spans" for individual record processing.
	for _, record := range records {
		// StartConsumer extracts the producer context from 'record.Context' and sets it as the PARENT.
		// We use 'originalCtx' to ensure these inner spans are siblings to the batch span, not children.
		_, innerSpan := StartConsumer(originalCtx, record, 1)

		// Create a link from the individual message trace back to the Batch Span.
		spanCtx := span.SpanContext()
		if spanCtx.IsValid() {
			innerSpan.AddLink(trace.Link{SpanContext: spanCtx})
		}
		innerSpan.End()
	}

	return ctx, span
}

func Start(ctx context.Context) (context.Context, trace.Span) {
	return tracer.Start(ctx, caller.FuncName(caller.WithSkip(1)))
}
