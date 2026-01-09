package observer

import (
	"fmt"
	"sync"
)

// EventHandler is a function that handles an event
type EventHandler func()

// EventBus manages event handlers
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

var globalEventBus = &EventBus{
	handlers: make(map[string][]EventHandler),
}

// Register registers an event handler
func Register(event string, handler EventHandler) {
	globalEventBus.mu.Lock()
	defer globalEventBus.mu.Unlock()
	globalEventBus.handlers[event] = append(globalEventBus.handlers[event], handler)
	fmt.Println("globalEventBus.handlers[event]:", event)
}

// UnregisterAll removes all handlers for an event
func UnregisterAll(event string) {
	globalEventBus.mu.Lock()
	defer globalEventBus.mu.Unlock()
	delete(globalEventBus.handlers, event)
}

// Trigger fires an event and calls all registered handlers
func Trigger(event string) {
	globalEventBus.mu.RLock()
	defer globalEventBus.mu.RUnlock()

	if handlers, ok := globalEventBus.handlers[event]; ok {
		for _, handler := range handlers {
			go handler() // async
		}
	}
}
