package event

import (
	ptypes "github.com/golang/protobuf/ptypes"
	uuid "github.com/google/uuid"
)

func NewEvent(source string) *Event {
  id := uuid.New()
  return &Event{
		Uuid: id[:],
		EventTime: ptypes.TimestampNow(),
		Source: source,
	}
}

func (e Event) DoProcessing() {
	e.Nprocessed++
}

func (e Event) GetId() string {
	id, _ := uuid.FromBytes(e.GetUuid())
	return id.String()
}
