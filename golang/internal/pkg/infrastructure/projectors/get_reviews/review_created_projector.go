package get_reviews

import (
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/queries"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/mysql"
)

func NewReviewCreatedProjector(getReviewsRepository *mysql.QueriesReviewsRepository) *ReviewCreatedProjector {
	return &ReviewCreatedProjector{getReviewsRepository: getReviewsRepository}
}

type ReviewCreatedProjector struct {
	getReviewsRepository *mysql.QueriesReviewsRepository
}

func (p *ReviewCreatedProjector) Project(review *commands.ReviewCreatedEvent) error {
	logrus.Info("Projecting")

	return p.getReviewsRepository.AddOrUpdateByUUID(
		&queries.Review{
			UUID:    review.ReviewUUID,
			Comment: review.Comment,
			Rating:  review.Rating,
		},
	)
}
