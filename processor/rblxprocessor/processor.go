package rblxprocessor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
	"strconv"
)

type rblxProcessor struct {
	attachedKey string
	attachedValue string
	logger          *zap.Logger
}

const (
	sourceFormat = "rblx"
)

func newRblxProcessor(logger *zap.Logger, nextConsumer consumer.Traces, cfg Config) (*rblxProcessor, error) {
	if nextConsumer == nil {
		return nil, component.ErrNilNextConsumer
	}

	rblx := &rblxProcessor{attachedKey: cfg.AttachedKey, attachedValue: cfg.AttachedValue, logger: logger}
	logger.Info("rblx processor is built")
	return rblx, nil
}



func (rblx *rblxProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {

	rblx.logger.Error("rblx processing trace, resources length" + strconv.Itoa(td.ResourceSpans().Len()))
	rss := td.ResourceSpans()

	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		ilss := rs.ScopeSpans()
		_ = rs.Resource()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			spans := ils.Spans()
			_ = ils.Scope()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)

				rblx.logger.Info("rblx span name is " + span.Name())
				span.Attributes().PutStr(rblx.attachedKey,rblx.attachedValue)
			}
		}
	}
	return td, nil
}

func (rblx *rblxProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (rblx *rblxProcessor) Start(_ context.Context, host component.Host) error {
	return nil
}

func (rblx *rblxProcessor) Shutdown(context.Context) error {
	return nil
}

