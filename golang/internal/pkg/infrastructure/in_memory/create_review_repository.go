package in_memory

import (
	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

type CreateReviewRepository struct {
	reviews []*domain.Review
}

func (r *CreateReviewRepository) Add(review *domain.Review) error {
	r.reviews = append(r.reviews, review)

	return nil
}
