package bus

type EventMeta struct {
	Type     string `json:"type"`
	Resource string `json:"resource"`
	Id       string `json:"id,omitempty"`
}

type Event[T interface{}] struct {
	Metadata EventMeta `json:"metadata"`
	Data     T         `json:"data"`
}
