package envelope

import (
	"encoding/json"
	"time"
)

type (
	// Envelope はイベントストアに保存するイベントの構造体
	Envelope struct {
		AggregateID string          `json:"aggregate_id"`
		SeqNr       uint64          `json:"seq_nr"`
		EventType   string          `json:"event_type"`
		Payload     json.RawMessage `json:"payload"`
		OccurredAt  time.Time       `json:"occurred_at"`
	}
)

func NewEnvelope(
	aggregateID string,
	seqNr uint64,
	eventType string,
	payload interface{},
) (*Envelope, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Envelope{
		AggregateID: aggregateID,
		SeqNr:       seqNr,
		EventType:   eventType,
		Payload:     b,
	}, nil
}
