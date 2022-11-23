package bus

import "encoding/json"

// EventMeta defines a standardized set of properties that all broadcase events
// must provide for identification.
type EventMeta struct {
	// Type is the type of event that has been emitted.
	Type string `json:"type"`

	// Resource is the resource type this event is relevant to (i.e. `class`).
	Resource string `json:"resource"`

	// Id is the resource ID this event is relevant to, but is optional and should
	// only be included with some event types (i.e. create, update, etc.).
	Id string `json:"id,omitempty"`
}

// Event is the standard OpenSchool event bus message format, providing metadata
// for each event that gets broadcast.
type Event[T interface{}] struct {
	// Metadata is the standard [EventMeta] information struct for the event.
	Metadata EventMeta `json:"metadata"`

	// Data is the payload of the event, which should be JSON-serializable.
	Data T `json:"data"`
}

// String returns a the JSON string form of the [Event], panicking if the Data
// field is not JSON-serializable for any reason.
func (e Event[T]) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(b)
}
