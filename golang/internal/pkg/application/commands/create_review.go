package commands

import (
	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

func NewCreateReviewCommandHandler(
	reviewRepository domain.ReviewRepository,
) *CreateReviewCommandHandler {
	return &CreateReviewCommandHandler{
		reviewRepository: reviewRepository,
	}
}

type CreateReviewCommandHandler struct {
	reviewRepository domain.ReviewRepository
	eventRepository  EventRepository
	transactionFactory TransactionFactory
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

type EventRepository interface {
	Save(name string, payload interface{}) error
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

	err := h.reviewRepository.Save(review)
	if err != nil {
		return err
	}

	err = h.eventRepository.Save(
		"review_created",
		ReviewCreatedEvent{
			ReviewUUID: review.Uuid().String(),
			Comment:    review.Comment(),
			Rating:     review.Rating(),
		},
	)

	return nil
}
