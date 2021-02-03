package streamer

import (
	"github.com/oleglarionov/otusgolang_finalproject/internal/application/event"
)

type AMQPStreamer struct {
}

func NewAMQPStreamer() *AMQPStreamer {
	return &AMQPStreamer{}
}

func (s *AMQPStreamer) Push(event event.Event) error {
	return nil // todo
}
