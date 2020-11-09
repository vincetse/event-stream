package main

import (
	log "github.com/sirupsen/logrus"
	event "github.com/vincetse/event-stream/pkg/event/v1"
)

type EventGenerator struct {
	// name of the event generator
	Name string
}

func NewEventGenerator(name string) *EventGenerator {
	gen := &EventGenerator{
		Name: name,
	}
	log.Infof("setting generator name to %s", name)
	return gen
}

func (gen *EventGenerator) NewEvent() *event.Event {
	name := gen.Name
	e := event.NewEvent(name)
	return e
}
