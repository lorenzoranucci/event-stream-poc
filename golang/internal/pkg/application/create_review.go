package application

import (
	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

type CreateReviewCommandHandler struct {
	reviewRepository CreateReviewRepository
	eventBus         EventBus
}

func NewCreateReviewCommandHandler(reviewRepository CreateReviewRepository, eventBus EventBus) *CreateReviewCommandHandler {
	return &CreateReviewCommandHandler{reviewRepository: reviewRepository, eventBus: eventBus}
}

type CreateReviewCommand struct {
	Comment string
	Rating  int32
}

type CreateReviewRepository interface {
	Add(review *domain.Review) error
}

type EventBus interface {
	DispatchEvent(event ReviewCreatedEvent)
}

type ReviewCreatedEvent struct {
	UUID    string
	Comment string
	Rating  int32
}

func (h *CreateReviewCommandHandler) Execute(command CreateReviewCommand) error {
	review := domain.NewReview(command.Comment, command.Rating)

	err := h.reviewRepository.Add(review)
	if err != nil {
		return err
	}

	h.eventBus.DispatchEvent(
		ReviewCreatedEvent{
			UUID:    review.Uuid().String(),
			Comment: review.Comment(),
			Rating:  review.Rating(),
		},
	)

	return nil
}
