package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"
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
	event commands.IntegrationEvent,
	transaction commands.Transaction,
) error {
	eventPayload, ok := event.Payload.(commands.ReviewCreatedEvent)
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

	lastMessageCounterByAggregate, err := selectLastMessageCounterByAggregate(event.AggregateID, i.db)
	if err != nil {
		return err
	}

	// this approach assumes that (aggregate_id, message_counter_by_aggregate) are unique and
	// the Read Uncommitted isolation level is not used.
	// In this way the counter will be perfectly consequential for each aggregate.
	messageCounterByAggregate := lastMessageCounterByAggregate + 1

	_, err = tx.Exec(
		"INSERT INTO review_events_outbox "+
			"(uuid, aggregate_id, name, payload, version, status, message_counter_by_aggregate) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE status=?",
		event.UUID,
		event.AggregateID,
		event.Name,
		payload,
		event.Version,
		int(event.Status),
		messageCounterByAggregate,
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

func selectLastMessageCounterByAggregate(aggregateId string, db *sqlx.DB) (int32, error) {
	rows, err := db.Query(
		`SELECT MAX(message_counter_by_aggregate)
			FROM review_events_outbox 
			WHERE aggregate_id = ?`,
		aggregateId,
	)
	if err != nil {
		return 0, err
	}

	rows.Next()
	var lastMessageCounterByAggregateId sql.NullInt32
	err = rows.Scan(&lastMessageCounterByAggregateId)
	if err != nil {
		return 0, err
	}

	var result int32
	if lastMessageCounterByAggregateId.Valid {
		result = lastMessageCounterByAggregateId.Int32
	}

	return result, err
}
