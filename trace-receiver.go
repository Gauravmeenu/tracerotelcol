package tracerotelcol

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type tracerotelcolReceiver struct{
	host component.Host
	cancel context.CancelFunc
	logger *zap.Logger
	nextConsumer consumer.Traces
	config *Config

}

func (tracerotelcolRcvr *tracerotelcolReceiver)Start(ctx context.Context, host component.Host)error{
	tracerotelcolRcvr.host=host
	ctx=context.Background()
	ctx, tracerotelcolRcvr.cancel=context.WithCancel(ctx)
	interval, _ := time.ParseDuration(tracerotelcolRcvr.config.Interval)
	go func(){
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for{
			select{
			case <-ticker.C:
				tracerotelcolRcvr.logger.Info("I should start processing now!!")
			    tracerotelcolRcvr.nextConsumer.ConsumeTraces(ctx, generateTraces())
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (tracerotelcolRcvr *tracerotelcolReceiver)Shutdown(ctx context.Context)error{
	tracerotelcolRcvr.cancel()
	return nil
}