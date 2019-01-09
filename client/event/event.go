package event

// Events channel.
type Events chan *Event

// Emit an event.
func (e Events) Emit(event Event) {
	e <- &event
}

// Event is a representation of an operation performed
type Event struct {
	Name  string
	Value string
}
