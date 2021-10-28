package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

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

func (i *IntegrationReviewEventsOutboxRepository) Save(
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

	tx, shouldCommit, err := getTransaction(transaction, i.db)
	if err != nil {
		return err
	}

	if shouldCommit {
		defer func(tx *sql.Tx) {
			err := tx.Rollback()
			if err != nil {
				logrus.Error(err)
			}
		}(tx)
	}

	_, err = tx.Exec(
		"INSERT INTO review_events_outbox (uuid, aggregate_id, name, payload, version, status) "+
			"VALUES (?, ?, ?, ?, ?,?)"+
			"ON DUPLICATE KEY UPDATE aggregate_id=?, name=?, payload=?, version=?, status=?",
		event.UUID,
		event.AggregateID,
		event.Name,
		payload,
		event.Version,
		int(event.Status),
		event.AggregateID,
		event.Name,
		payload,
		event.Version,
		int(event.Status),
	)
	if err != nil {
		return err
	}

	if shouldCommit {
		return tx.Commit()
	}

	return nil
}
