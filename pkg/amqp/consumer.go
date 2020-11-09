package amqp

import (
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	mq "github.com/streadway/amqp"
	event "github.com/vincetse/event-stream/pkg/event/v1"
	opts "github.com/vincetse/event-stream/pkg/flags"
)

type Consumer struct {
	options *opts.ConsumerFlags
	conn *mq.Connection
	ch *mq.Channel
	queue mq.Queue
	done    chan error
}

func NewConsumer(o *opts.ConsumerFlags) *Consumer {
	return &Consumer{
		options: o,
		done: make(chan error),
	}
}

func (p *Consumer) Open() (err error) {
	p.conn = Dial(p.options.Uri)
	ch, err := p.conn.Channel()
	if err == nil {
		log.Infof("got channel")
		p.ch = ch
	} else {
		log.Debug(err)
		return err
	}

	// use the default exchange if none is give in the command-line
	// parameter.
	if p.options.ExchangeName != "" {
		log.Infof("Declaring exchange %s", p.options.ExchangeName)
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
			log.Debug(err)
			return err
		}
	} else {
		log.Debugf("Using default exchange")
	}

	log.Infof("Declaring queue %s", p.options.QueueName)
	queue, err := ch.QueueDeclare(
		p.options.QueueName,
		false, // durable?
		false, // auto-deleted?
		false, // exclusive?
		false, // no wait?
		nil,
	)
	if err != nil {
		log.Debugf("Error declaring queue %s", queue.Name)
		log.Debug(err)
		return err
	}
	p.queue = queue

	if p.options.ExchangeName != "" {
		log.Infof("Binding to queue %s", p.options.QueueName)
		if err = ch.QueueBind(
			queue.Name,
			p.options.RoutingKey,
			p.options.ExchangeName,
			false, // no wait?
			nil,
		); err != nil {
			log.Debugf("Error binding queue %s", p.options.QueueName)
			log.Debug(err)
			return err
		}
	}

	return err
}

func (p *Consumer) Close() {
	if p.conn != nil {
		p.conn.Close()
		log.Infof("closed connection to %s", p.options.Uri)
	}
}

func (p *Consumer) Consume() {
	ch := p.ch

	deliveries, err := ch.Consume(
		p.queue.Name,
		"",
		false, // no ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,
	)
	if err != nil {
		log.Debug(err)
		return
	}

	go handle(p.options.RoutingKey, deliveries, p.done)
}

func handle(routing_key string, deliveries <-chan mq.Delivery, done chan error) {
	for d := range deliveries {
		e := &event.Event{}
		if err := proto.Unmarshal(d.Body, e); err != nil {
			log.Error(err)
		} else {
			log.Debugf("consumed event %s [routing-key=%-32s] [source-routing-key=%-32s]", e.GetId(), routing_key, e.GetRoutingKey())
		}

		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
