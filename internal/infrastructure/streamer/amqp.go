package streamer

import (
	"github.com/oleglarionov/otusgolang_finalproject/internal/application/event"
)

type AmqpStreamer struct {
}

func NewAmqpStreamer() *AmqpStreamer {
	return &AmqpStreamer{}
}

func (s *AmqpStreamer) Push(event event.Event) error {
	return nil // todo
}
