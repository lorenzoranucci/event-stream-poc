package in_memory

import (
	"math"

	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
)

type GetReviewsRepository struct {
	reviews []application.Review
}

func (r *GetReviewsRepository) Find(query application.GetReviewsQuery) ([]application.Review, error) {
	elementsAfterOffset := int64(len(r.reviews)) - query.Offset

	logrus.Infof(
		"finding limit: %d offset: %d, len: %d, aft: %d",
		query.Limit,
		query.Offset,
		len(r.reviews),
		elementsAfterOffset,
	)

	if elementsAfterOffset <= 0 {
		logrus.Infof("not enough elements")

		return []application.Review{}, nil
	}

	min := math.Min(float64(elementsAfterOffset), float64(query.Limit))

	return r.reviews[query.Offset:int(min)], nil
}

func (r *GetReviewsRepository) Add(review *application.Review) error {
	r.reviews = append(r.reviews, *review)

	return nil
}
