package event

type Event struct {
	ID            int64       `json:"id"`
	AggregateID   string      `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	Version       int         `json:"version"`
	EventName     string      `json:"event_name"`
	EventData     interface{} `json:"event_data"`
	CreatedAt     int64       `json:"created_at"`
}
