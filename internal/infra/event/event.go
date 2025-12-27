package event

import (
	"context"
	"reflect"
	"sync"
)

// Event is a marker interface for events
type Event interface {
	EventName() string
}

// Listener handles an event
type Listener interface {
	Handle(ctx context.Context, event Event) error
}

// ListenerFunc is a function adapter for Listener
type ListenerFunc func(ctx context.Context, event Event) error

func (f ListenerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

// Dispatcher manages event dispatching
type Dispatcher struct {
	mu        sync.RWMutex
	listeners map[string][]Listener
	async     bool
}

// dispatcher is the global dispatcher instance
var (
	globalDispatcher *Dispatcher
	once             sync.Once
)

// Global returns the global dispatcher instance
func Global() *Dispatcher {
	once.Do(func() {
		globalDispatcher = New()
	})
	return globalDispatcher
}

// New creates a new event dispatcher
func New() *Dispatcher {
	return &Dispatcher{
		listeners: make(map[string][]Listener),
	}
}

// Listen registers a listener for an event
func (d *Dispatcher) Listen(eventName string, listener Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listeners[eventName] = append(d.listeners[eventName], listener)
}

// ListenFunc registers a function listener for an event
func (d *Dispatcher) ListenFunc(eventName string, fn func(ctx context.Context, event Event) error) {
	d.Listen(eventName, ListenerFunc(fn))
}

// Subscribe registers a listener for a typed event using reflection
func (d *Dispatcher) Subscribe(eventType Event, listener Listener) {
	d.Listen(eventType.EventName(), listener)
}

// Dispatch fires an event to all registered listeners
func (d *Dispatcher) Dispatch(ctx context.Context, event Event) error {
	d.mu.RLock()
	listeners := d.listeners[event.EventName()]
	d.mu.RUnlock()

	for _, listener := range listeners {
		if err := listener.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// DispatchAsync fires an event asynchronously to all registered listeners
func (d *Dispatcher) DispatchAsync(ctx context.Context, event Event) {
	d.mu.RLock()
	listeners := d.listeners[event.EventName()]
	d.mu.RUnlock()

	for _, listener := range listeners {
		go func(l Listener) {
			_ = l.Handle(ctx, event)
		}(listener)
	}
}

// DispatchAll fires multiple events
func (d *Dispatcher) DispatchAll(ctx context.Context, events ...Event) error {
	for _, event := range events {
		if err := d.Dispatch(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// HasListeners checks if an event has any listeners
func (d *Dispatcher) HasListeners(eventName string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.listeners[eventName]) > 0
}

// Forget removes all listeners for an event
func (d *Dispatcher) Forget(eventName string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.listeners, eventName)
}

// Flush removes all listeners
func (d *Dispatcher) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listeners = make(map[string][]Listener)
}

// GetEventName extracts event name from type (helper for struct events)
func GetEventName(event Event) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath() + "." + t.Name()
}

// --- Convenience functions using global dispatcher ---

// Listen registers a listener on the global dispatcher
func Listen(eventName string, listener Listener) {
	Global().Listen(eventName, listener)
}

// ListenFunc registers a function listener on the global dispatcher
func ListenFunc(eventName string, fn func(ctx context.Context, event Event) error) {
	Global().ListenFunc(eventName, fn)
}

// Dispatch fires an event on the global dispatcher
func Dispatch(ctx context.Context, event Event) error {
	return Global().Dispatch(ctx, event)
}

// DispatchAsync fires an event asynchronously on the global dispatcher
func DispatchAsync(ctx context.Context, event Event) {
	Global().DispatchAsync(ctx, event)
}
