package in_memory

import (
	"math"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

type GetReviewsRepository struct {
	reviews []domain.Review
}

func (g *GetReviewsRepository) Find(query application.GetReviewsQuery) ([]domain.Review, error) {
	elementsAfterOffset := int64(len(g.reviews)) - query.Offset
	if elementsAfterOffset <= 0 {
		return []domain.Review{}, nil
	}

	min := math.Min(float64(elementsAfterOffset), float64(query.Limit))

	return g.reviews[query.Offset:int(min)], nil
}
