package commands

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

const (
	reviewCreatedEventVersion = "0.1.0"
	reviewCreatedEventName    = "review_created"
)

func NewCreateReviewCommandHandler(
	reviewRepository ReviewRepository,
	transactionManager TransactionManager,
	eventOutboxRepository IntegrationEventOutboxRepository,
	eventBus IntegrationEventBus,
) *CreateReviewCommandHandler {
	return &CreateReviewCommandHandler{
		reviewRepository:      reviewRepository,
		transactionManager:    transactionManager,
		eventOutboxRepository: eventOutboxRepository,
		eventBus:              eventBus,
	}
}

type CreateReviewCommandHandler struct {
	reviewRepository      ReviewRepository
	transactionManager    TransactionManager
	eventOutboxRepository IntegrationEventOutboxRepository
	eventBus              IntegrationEventBus
}

type CreateReviewCommand struct {
	Comment string
	Rating  int32
}

type ReviewCreatedEvent struct {
	ReviewUUID string
	Comment    string
	Rating     int32
}

func (h *CreateReviewCommandHandler) Execute(command CreateReviewCommand) error {
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

func (h *CreateReviewCommandHandler) executeTransactionally(
	command CreateReviewCommand,
	transaction Transaction,
) ([]IntegrationEvent, error) {
	review := domain.NewReview(command.Comment, command.Rating)

	err := h.reviewRepository.Save(review, transaction)
	if err != nil {
		return nil, err
	}

	event := IntegrationEvent{
		UUID:        uuid.New().String(),
		AggregateID: review.Uuid().String(),
		Name:        reviewCreatedEventName,
		Payload: ReviewCreatedEvent{
			ReviewUUID: review.Uuid().String(),
			Comment:    review.Comment(),
			Rating:     review.Rating(),
		},
		Status:  ToBeDispatched,
		Version: reviewCreatedEventVersion,
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

func (h *CreateReviewCommandHandler) executePostTransaction(events []IntegrationEvent) error {
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
