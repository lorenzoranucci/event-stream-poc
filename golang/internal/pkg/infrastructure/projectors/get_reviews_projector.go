package projectors

import (
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/mysql"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
)

func NewGetReviewsProjector(getReviewsRepository *mysql.GetReviewsRepository) *GetReviewsProjector {
	return &GetReviewsProjector{getReviewsRepository: getReviewsRepository}
}

type GetReviewsProjector struct {
	getReviewsRepository *mysql.GetReviewsRepository
}

func (p *GetReviewsProjector) Project(review *application.ReviewCreatedEvent) error {
	logrus.Info("Projecting")

	return p.getReviewsRepository.AddOrUpdateByUUID(
		&application.Review{
			UUID:    review.UUID,
			Comment: review.Comment,
			Rating:  review.Rating,
		},
	)
}
