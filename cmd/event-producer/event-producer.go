package main

import (
	"os"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/vincetse/event-stream/pkg/amqp"
	myflags "github.com/vincetse/event-stream/pkg/flags"
)

func main() {
	log.SetLevel(log.DebugLevel)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	options := myflags.NewProducerFlags("")
	myflags.Parse()

	q := amqp.NewProducer(options)
	err = q.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()

	gen := NewEventGenerator(hostname)

	var i int64
	for i = 0; i < options.EventCount; i++ {
		e := gen.NewEvent()
		q.Publish(e)
		// pause so we don't hog the CPU
		time.Sleep(0.5 * 1e9)
	}
}
