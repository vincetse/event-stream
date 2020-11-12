package flags

import (
	"github.com/spf13/pflag"
)

type ProducerFlags struct {
	Uri string
	ExchangeName string
	ExchangeType string
	RoutingKey string
	EventCount int64
}

func setupProducerFlags(f *ProducerFlags, prefix string) {
	pflag.StringVar(&f.Uri, prefix + "uri", "", "Producer AMPQ URI")
	pflag.StringVar(&f.ExchangeName, prefix + "exchange-name", "", "Producer AMPQ exchange name")
	pflag.StringVar(&f.ExchangeType, prefix + "exchange-type", "direct", "Producer AMPQ exchange type")
	pflag.StringVar(&f.RoutingKey, prefix + "routing-key", "", "Producer AMPQ routing-key")
	pflag.Int64Var(&f.EventCount, prefix + "event-count", 0, "Event count (no effect for processors")
}

func NewProducerFlags(prefix string) *ProducerFlags {
	retv := &ProducerFlags{}
	setupProducerFlags(retv, prefix)
	return retv
}
