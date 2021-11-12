package commands

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	reviewRatingIncrementEventEventVersion = "0.1.0"
	reviewRatingIncrementEventName         = "review_rating_incremented"
)

func NewIncrementReviewRatingCommandHandler(
	reviewRepository ReviewRepository,
	transactionManager TransactionManager,
	eventOutboxRepository IntegrationEventOutboxRepository,
	eventBus IntegrationEventBus,
) *IncrementReviewRatingCommandHandler {
	return &IncrementReviewRatingCommandHandler{
		reviewRepository:      reviewRepository,
		transactionManager:    transactionManager,
		eventOutboxRepository: eventOutboxRepository,
		eventBus:              eventBus,
	}
}

type IncrementReviewRatingCommandHandler struct {
	reviewRepository      ReviewRepository
	transactionManager    TransactionManager
	eventOutboxRepository IntegrationEventOutboxRepository
	eventBus              IntegrationEventBus
}

type IncrementReviewRatingCommand struct {
	ReviewUUID uuid.UUID
}

type ReviewRatingIncrementedEvent struct {
	ReviewUUID string
}

func (h *IncrementReviewRatingCommandHandler) Execute(command IncrementReviewRatingCommand) error {
	transaction, err := openTransaction(h.transactionManager)
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	events, err := h.executeTransactionally(command, transaction)
	if err != nil {
		return err
	}

	return h.executePostTransaction(events)
}

func (h *IncrementReviewRatingCommandHandler) executeTransactionally(
	command IncrementReviewRatingCommand,
	transaction Transaction,
) ([]IntegrationEvent, error) {
	review, err := h.reviewRepository.FindByUUID(command.ReviewUUID)
	if err != nil {
		return nil, err
	}

	if review == nil {
		return nil, fmt.Errorf("review not found")
	}

	review.IncrementRating()

	err = h.reviewRepository.Save(review, transaction)
	if err != nil {
		return nil, err
	}

	event := IntegrationEvent{
		UUID:        uuid.New().String(),
		AggregateID: review.Uuid().String(),
		Name:        reviewRatingIncrementEventName,
		Payload: ReviewRatingIncrementedEvent{
			ReviewUUID: review.Uuid().String(),
		},
		Status:  ToBeDispatched,
		Version: reviewRatingIncrementEventEventVersion,
	}

	err = h.eventOutboxRepository.Save(event, transaction)
	if err != nil {
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return []IntegrationEvent{event}, nil
}

func (h *IncrementReviewRatingCommandHandler) executePostTransaction(events []IntegrationEvent) error {
	for _, event := range events {
		h.eventBus.DispatchEvent(
			event,
		)

		event.Status = Dispatched
		// error is not handle here since it's acceptable to have the event dispatched again (at least once).
		// The event may be dispatched again by a dedicated job.
		err := h.eventOutboxRepository.Save(event, nil)
		if err != nil {
			logrus.Error(err)
		}
	}

	return nil
}
