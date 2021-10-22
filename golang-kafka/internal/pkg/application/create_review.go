package application

import (
	"github.com/ProntoPro/golang-kafka/internal/pkg/domain"
)

type CreateReviewCommandHandler struct {
	reviewRepository ReviewRepository
	eventBus          EventBus
}

type CreateReviewCommand struct {
	Comment string
	Rating int
}

type ReviewRepository interface {
	Add(tu *domain.Review) error
}

type EventBus interface {
	DispatchEvent(event ReviewCreatedEvent)
}

type ReviewCreatedEvent struct {
	Review *domain.Review
}

func (h *CreateReviewCommandHandler) Execute(command CreateReviewCommand) error {
	review := domain.NewReview(command.Comment, command.Rating)

	err := h.reviewRepository.Add(review)
	if err != nil {
		return err
	}

	h.eventBus.DispatchEvent(ReviewCreatedEvent{Review: review})

	return nil
}
