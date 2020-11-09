package amqp

import (
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	mq "github.com/streadway/amqp"
	event "github.com/vincetse/event-stream/pkg/event/v1"
	opts "github.com/vincetse/event-stream/pkg/flags"
)

type Producer struct {
	options *opts.ProducerFlags
	conn *mq.Connection
	ch *mq.Channel
}

func NewProducer(o *opts.ProducerFlags) *Producer {
	return &Producer{
		options: o,
	}
}

func (p *Producer) Open() (err error) {
	p.conn = Dial(p.options.Uri)
	ch, err := p.conn.Channel()
	if err == nil {
		log.Infof("got channel")
		p.ch = ch
	}
	// TODO handle error

	// use the default exchange if none is give in the command-line
	// parameter.
	if p.options.ExchangeName != "" {
		if err := ch.ExchangeDeclare(
			p.options.ExchangeName,
			p.options.ExchangeType,
			false, // durable?
			true, // auto-deleted?
			false, // internal?
			false, // no wait?
			nil,
		); err != nil {
			log.Debugf("Error declaring exchange %s/%s",
				p.options.ExchangeName,
				p.options.ExchangeType,
			)
		}
	} else {
		log.Debugf("Using default exchange")
	}

	return err
}

func (p *Producer) Close() {
	if p.conn != nil {
		p.conn.Close()
		log.Infof("closed connection to %s", p.options.Uri)
	}
}

func (p *Producer) Publish(e *event.Event) {
	e.RoutingKey = p.options.RoutingKey
	data, err := proto.Marshal(e)
	ch := p.ch

	err = ch.Publish(
		p.options.ExchangeName,
		p.options.RoutingKey,
		false,
		false,
		mq.Publishing {
			ContentType: "text/plain",
			Body: data,
		},
	)

	if err != nil {
		log.Error(err)
	} else {
		log.Debugf("publishing event %s [routing-key=%-32s]", e.GetId(), p.options.RoutingKey)
	}
}
