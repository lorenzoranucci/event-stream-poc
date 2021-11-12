package commands

import (
	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
	"github.com/google/uuid"
)

func NewCreateReviewCommandHandler(
	reviewRepository ReviewRepository,
) *CreateReviewCommandHandler {
	return &CreateReviewCommandHandler{
		reviewRepository: reviewRepository,
	}
}

type CreateReviewCommandHandler struct {
	reviewRepository ReviewRepository
	eventRepository  EventRepository
	transactionFactory TransactionFactory
}

type CreateReviewCommand struct {
	Comment string
	Rating  int32
}

type ReviewRepository interface {
	SaveTransactional(transaction Transaction, review *domain.Review) error
	FindByUUID(reviewUUID uuid.UUID) (*domain.Review, error)
}

type ReviewCreatedEvent struct {
	ReviewUUID string
	Comment    string
	Rating     int32
}

type EventRepository interface {
	SaveTransactional(transaction Transaction, name string, payload interface{}) error
}

type Transaction interface {
	Begin() error
	Commit() error
}

type TransactionFactory interface {
	Create() (Transaction, error)
}

func (h *CreateReviewCommandHandler) Execute(command CreateReviewCommand) error {
	review := domain.NewReview(command.Comment, command.Rating)

	transaction, err := h.transactionFactory.Create()
	if err != nil {
		return err
	}

	err = transaction.Begin()
	if err != nil {
		return err
	}

	err = h.reviewRepository.SaveTransactional(transaction, review)
	if err != nil {
		return err
	}

	err = h.eventRepository.SaveTransactional(
		transaction,
		"review_created",
		ReviewCreatedEvent{
			ReviewUUID: review.Uuid().String(),
			Comment:    review.Comment(),
			Rating:     review.Rating(),
		},
	)

	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}
