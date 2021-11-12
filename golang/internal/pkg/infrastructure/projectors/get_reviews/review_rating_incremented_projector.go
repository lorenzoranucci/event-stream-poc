package get_reviews

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/mysql"
)

func NewReviewRatingIncrementedProjector(
	getReviewsRepository *mysql.QueriesReviewsRepository,
) *ReviewRatingIncrementedProjector {
	return &ReviewRatingIncrementedProjector{getReviewsRepository: getReviewsRepository}
}

type ReviewRatingIncrementedProjector struct {
	getReviewsRepository *mysql.QueriesReviewsRepository
}

func (p *ReviewRatingIncrementedProjector) Project(review *commands.ReviewRatingIncrementedEvent) error {
	logrus.Info("Projecting")

	uuidFromString, err := uuid.Parse(review.ReviewUUID)
	if err != nil {
		return err
	}

	storedReview, err := p.getReviewsRepository.FindByUUID(uuidFromString)
	if err != nil {
		return err
	}

	storedReview.Rating++

	return p.getReviewsRepository.AddOrUpdateByUUID(
		storedReview,
	)
}
