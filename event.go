package polldance

import (
	"log"
	"sync"
)

type EventHandler func(*EventData) error

type EventData struct {
	Source string
	Data   string
}
type EventProcessor struct {
	handlers []EventHandler
	wg       sync.WaitGroup
}

func (ec *EventProcessor) Push(e *EventData) {
	for _, h := range ec.handlers {
		if err := h(e); err != nil {
			log.Printf("error at handler: %s", err)
		}
	}
}

func (ec *EventProcessor) AddHandler(h EventHandler) {
	ec.handlers = append(ec.handlers, h)
}

func (ec *EventProcessor) Wait() {
	ec.wg.Wait()
}

func (ec *EventProcessor) AddSource() {
	ec.wg.Add(1)
}

func (ec *EventProcessor) RemoveSource() {
	ec.wg.Done()
}

func NewEventProcessor() *EventProcessor {
	return &EventProcessor{
		handlers: []EventHandler{},
	}
}
