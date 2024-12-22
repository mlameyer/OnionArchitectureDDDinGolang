package services

import "log"

type EventPublisher interface {
	Publish(event interface{}) error
}

type LoggerEventPublisher struct{}

func (p *LoggerEventPublisher) Publish(event interface{}) error {
	log.Printf("Event published: %+v", event)
	return nil
}
