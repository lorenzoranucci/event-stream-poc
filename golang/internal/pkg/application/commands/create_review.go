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
	review := domain.NewReview(command.Comment, command.Rating)

	return h.reviewRepository.Save(review)
}
