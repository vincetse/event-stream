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

	cOptions := myflags.NewConsumerFlags("consumer-")
	pOptions := myflags.NewProducerFlags("producer-")
	myflags.Parse()

	q := amqp.NewConsumer(cOptions)
	if err := q.Open(); err != nil {
		log.Fatal(err)
	}
	p := amqp.NewProducer(pOptions)
	x := NewProcessor(p)
	if err := x.Open(); err != nil {
		log.Fatal(err)
	}

	if err := q.Consume(x.Handle); err != nil {
		log.Fatal(err)
	}

	select {}
}

///////
type Processor struct {
	producer *amqp.Producer
}

func NewProcessor(producer *amqp.Producer) *Processor {
	return &Processor{
		producer: producer,
	}
}

func (p *Processor) Open() (error) {
	return p.producer.Open()
}

func (p *Processor) Handle(routing_key string, deliveries <-chan mq.Delivery, done chan error) {
	for d := range deliveries {
		e := &event.Event{}
		if err := proto.Unmarshal(d.Body, e); err != nil {
			log.Error(err)
		} else {
			e.DoProcessing()
			if err:= p.producer.Publish(e); err != nil {
				log.Errorf("Error processing event %s [routing-key=%-32s] [source-routing-key=%-32s]", e.GetId(), routing_key, e.GetRoutingKey())
			} else {
				log.Debugf("processed event %s [routing-key=%-32s] [source-routing-key=%-32s] [n=%d]", e.GetId(), routing_key, e.GetRoutingKey(), e.GetNprocessed())
			}
		}

		if err := d.Ack(false); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
