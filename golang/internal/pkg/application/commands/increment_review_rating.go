package commands

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

func NewIncrementReviewRatingCommandHandler(
	reviewRepository domain.ReviewRepository,
) *IncrementReviewRatingCommandHandler {
	return &IncrementReviewRatingCommandHandler{
		reviewRepository: reviewRepository,
	}
}

type IncrementReviewRatingCommandHandler struct {
	reviewRepository domain.ReviewRepository
}

type IncrementReviewRatingCommand struct {
	ReviewUUID uuid.UUID
}

type ReviewRatingIncrementedEvent struct {
	ReviewUUID string
}

func (h *IncrementReviewRatingCommandHandler) Execute(command IncrementReviewRatingCommand) error {
	review, err := h.reviewRepository.FindByUUID(command.ReviewUUID)
	if err != nil {
		return err
	}

	if review == nil {
		return fmt.Errorf("review not found")
	}

	review.IncrementRating()

	return h.reviewRepository.Save(review)
}
