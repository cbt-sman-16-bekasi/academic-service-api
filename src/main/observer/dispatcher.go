package observer

import (
	"fmt"
	"sync"
)

type EventHandler func()

type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

var globalEventBus = &EventBus{
	handlers: make(map[string][]EventHandler),
}

// Register event handler
func Register(event string, handler EventHandler) {
	globalEventBus.mu.Lock()
	defer globalEventBus.mu.Unlock()
	globalEventBus.handlers[event] = append(globalEventBus.handlers[event], handler)
	fmt.Println("globalEventBus.handlers[event]:", event)
}

// Unregister all handlers (optional)
func UnregisterAll(event string) {
	globalEventBus.mu.Lock()
	defer globalEventBus.mu.Unlock()
	delete(globalEventBus.handlers, event)
}

// Trigger the event
func Trigger(event string) {
	globalEventBus.mu.RLock()
	defer globalEventBus.mu.RUnlock()

	if handlers, ok := globalEventBus.handlers[event]; ok {
		for _, handler := range handlers {
			go handler() // async
		}
	}
}
