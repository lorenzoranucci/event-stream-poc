package application

import "github.com/ProntoPro/event-stream-golang/internal/pkg/domain"

type GetReviewsQueryHandler struct {
	reviewRepository GetReviewsRepository
}

func NewGetReviewsQueryHandler(reviewRepository GetReviewsRepository) *GetReviewsQueryHandler {
	return &GetReviewsQueryHandler{reviewRepository: reviewRepository}
}

type GetReviewsQuery struct {
	Limit  int32
	Offset int64
}

type GetReviewsRepository interface {
	Find(query GetReviewsQuery) ([]domain.Review, error)
}

func (h *GetReviewsQueryHandler) Execute(query GetReviewsQuery) ([]domain.Review, error) {
	return h.reviewRepository.Find(query)
}
