package rblxprocessor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	typeStr = "rblx"
	stability = component.StabilityLevelAlpha
)

func NewFactory() processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig,
		processor.WithTraces(createTraceProcessor, stability))
}

func createDefaultConfig() component.Config {
	return &Config{
		AttachedKey:   "testKey",
		AttachedValue: "testValue",
	}
}

func createTraceProcessor(
	ctx context.Context,
	set processor.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (processor.Traces, error) {
	rCfg := cfg.(*Config)
	rblx, err := newRblxProcessor(set.Logger, nextConsumer, *rCfg)
	if err != nil {
		return nil, err
	}

	return processorhelper.NewTracesProcessor(ctx, set, cfg, nextConsumer, rblx.processTraces, processorhelper.WithCapabilities(rblx.Capabilities()))
}
