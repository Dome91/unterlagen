package eventing

import (
	"errors"
	"github.com/rs/zerolog/log"
	"unterlagen/pkg/domain"
)

type SynchronousEventBus struct {
	handlers map[string][]domain.EventHandler
}

func NewEventBus() *SynchronousEventBus {
	return &SynchronousEventBus{
		handlers: make(map[string][]domain.EventHandler),
	}
}

func (b *SynchronousEventBus) Publish(event domain.Event) error {
	var errs error
	for _, handler := range b.handlers[event.Topic()] {
		err := handler.Handle(event)
		if err != nil {
			errors.Join(err, errs)
		}
	}

	if errs != nil {
		log.Err(errs).Msg("handling events failed")
	}

	return nil
}

func (b *SynchronousEventBus) Subscribe(topic string, handler domain.EventHandler) error {
	b.handlers[topic] = append(b.handlers[topic], handler)
	return nil
}
