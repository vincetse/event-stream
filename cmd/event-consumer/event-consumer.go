package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vincetse/event-stream/pkg/amqp"
	mq "github.com/streadway/amqp"
	event "github.com/vincetse/event-stream/pkg/event/v1"
	myflags "github.com/vincetse/event-stream/pkg/flags"
	"google.golang.org/protobuf/proto"
)

func main() {
	log.SetLevel(log.DebugLevel)

	options := myflags.NewConsumerFlags("")
	myflags.Parse()

	q := amqp.NewConsumer(options)
	if err := q.Open(); err != nil {
		log.Fatal(err)
	}
	if err := q.Consume(handle); err != nil {
		log.Fatal(err)
	}

	select {}
}

func handle(routing_key string, deliveries <-chan mq.Delivery, done chan error) {
	for d := range deliveries {
		e := &event.Event{}
		if err := proto.Unmarshal(d.Body, e); err != nil {
			log.Error(err)
		} else {
			log.Debugf("consumed event %s [routing-key=%-32s] [source-routing-key=%-32s]", e.GetId(), routing_key, e.GetRoutingKey())
		}

		if err := d.Ack(false); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
