package commands

import (
	"github.com/google/uuid"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

type ReviewRepository interface {
	Save(review *domain.Review, transaction Transaction) error
	FindByUUID(reviewUUID uuid.UUID) (*domain.Review, error)
}
