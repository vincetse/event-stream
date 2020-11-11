package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vincetse/event-stream/pkg/amqp"
	myflags "github.com/vincetse/event-stream/pkg/flags"
)

func main() {
	log.SetLevel(log.DebugLevel)

	options := myflags.NewConsumerFlags("")
	myflags.Parse()

	q := amqp.NewConsumer(options)
	if err := q.Open(); err != nil {
		log.Fatal(err)
	}
	if err := q.Consume(); err != nil {
		log.Fatal(err)
	}

	select {}
}
