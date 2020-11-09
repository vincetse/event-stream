package flags

import (
	"github.com/spf13/pflag"
)

type ConsumerFlags struct {
	Uri string
	ExchangeName string
	ExchangeType string
	RoutingKey string
	QueueName string
}

func setupConsumerFlags(f *ConsumerFlags, prefix string) {
	pflag.StringVar(&f.Uri, prefix + "uri", "", "Consumer AMPQ URI")
	pflag.StringVar(&f.ExchangeName, prefix + "exchange-name", "", "Consumer AMPQ exchange name")
	pflag.StringVar(&f.ExchangeType, prefix + "exchange-type", "direct", "Consumer AMPQ exchange type")
	pflag.StringVar(&f.RoutingKey, prefix + "routing-key", "", "Consumer AMPQ routing-key")
	pflag.StringVar(&f.QueueName, prefix + "queue-name", "", "Consumer AMPQ queue name")
}

func NewConsumerFlags(prefix string) *ConsumerFlags {
	retv := &ConsumerFlags{}
	setupConsumerFlags(retv, prefix)
	return retv
}
