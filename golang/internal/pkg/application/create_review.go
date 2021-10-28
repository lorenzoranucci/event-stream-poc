package application

import (
	"github.com/google/uuid"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

const (
	eventVersion = "0.1.0"
	eventName    = "review_created_event"
)

func NewCreateReviewCommandHandler(
	reviewRepository CreateReviewRepository,
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
	reviewRepository      CreateReviewRepository
	transactionManager    TransactionManager
	eventOutboxRepository IntegrationEventOutboxRepository
	eventBus              IntegrationEventBus
}

type CreateReviewCommand struct {
	Comment string
	Rating  int32
}

type CreateReviewRepository interface {
	Add(review *domain.Review, transaction Transaction) error
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

	err := h.reviewRepository.Add(review, transaction)
	if err != nil {
		return nil, err
	}

	event := IntegrationEvent{
		UUID:        uuid.New().String(),
		AggregateID: review.Uuid().String(),
		Name:        eventName,
		Payload: ReviewCreatedEvent{
			ReviewUUID: review.Uuid().String(),
			Comment:    review.Comment(),
			Rating:     review.Rating(),
		},
		Version: eventVersion,
	}

	err = h.eventOutboxRepository.Add(event, transaction)
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
	}

	return nil
}
