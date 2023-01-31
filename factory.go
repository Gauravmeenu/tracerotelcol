package tracerotelcol

import (
	"context"
	"strconv"

	"go.opentelemetry.io/collector/component"
	// "go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	typeStr         = "tracerotelcol"
	defaultInterval = 1
)

func CreateDefaultConfig() component.Config {

	return &Config{
		// ReceiverSettings: config.NewReceiverSettings(component.NewID(typeStr)),
		Interval:         strconv.Itoa(defaultInterval),
	}
}

// func CreateTracesReceiver(_ context.Context, params component.ReceiverCreateSettings, bcfg Config,
//
//	nextConsumer Component.consumer.Traces) (component.TracesReceiver, error){
//		return nil, nil
//	}
func createTracesReceiver(_ context.Context, p receiver.CreateSettings, cfg component.Config, consumer consumer.Traces) (receiver.Traces, error) {
	if consumer == nil{
		return nil, component.ErrNilNextConsumer
	}
	logger := p.Logger
	tracerotelcolCfg := cfg.(*Config)

	traceRcvr := &tracerotelcolReceiver{
		logger: logger,
		nextConsumer: consumer,
		config: tracerotelcolCfg,
	}

	return traceRcvr, nil
}

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		CreateDefaultConfig,
		receiver.WithTraces(createTracesReceiver, component.StabilityLevelStable),
	)

}