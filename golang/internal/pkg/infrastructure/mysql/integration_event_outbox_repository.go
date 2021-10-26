package mysql

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
)

type IntegrationReviewEventsOutboxRepository struct {
	db *sqlx.DB
}

func NewIntegrationReviewEventsOutboxRepository(db *sqlx.DB) *IntegrationReviewEventsOutboxRepository {
	return &IntegrationReviewEventsOutboxRepository{db: db}
}

type Payload struct {
	UUID    string `json:"uuid"`
	Comment string `json:"comment"`
	Rating  int32  `json:"rating"`
}

func (i *IntegrationReviewEventsOutboxRepository) Add(
	event application.IntegrationEvent,
	transaction application.Transaction,
) error {
	eventPayload, ok := event.Payload.(application.ReviewCreatedEvent)
	if !ok {
		return fmt.Errorf("unsupported event payload")
	}

	payload, err := json.Marshal(
		Payload{
			UUID:    eventPayload.ReviewUUID,
			Comment: eventPayload.Comment,
			Rating:  eventPayload.Rating,
		},
	)
	if err != nil {
		return err
	}

	tx, err := getTransaction(transaction, i.db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO review_events_outbox (uuid, aggregate_id, name, payload, version) "+
			"VALUES (?, ?, ?, ?, ?)",
		event.UUID,
		event.AggregateID,
		event.Name,
		payload,
		event.Version,
	)
	if err != nil {
		return err
	}

	return nil
}
