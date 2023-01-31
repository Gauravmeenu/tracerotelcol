package tracerotelcol

import (
	"fmt"
	"time"

	// "go.opentelemetry.io/collector/config"
)

type Config struct{
	// config.ReceiverSettings `mapstructure:",squash"`
	Interval string `mapstructure:"interval"`
	NumberOfTraces int `mapstructure:"number_of_traces"`
}

func (cfg *Config) Validate() error {
	interval, _ := time.ParseDuration(cfg.Interval)
	if (interval.Minutes()<1){
		return fmt.Errorf("interval must be at least one minute")
	}
	if (cfg.NumberOfTraces<1){
		return fmt.Errorf("number of traces should be greater than or equal to 1")
	}
	return nil
}