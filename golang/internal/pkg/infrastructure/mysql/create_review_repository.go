package mysql

import (
	"github.com/jmoiron/sqlx"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

func NewCreateReviewRepository(conn *sqlx.DB) *CreateReviewRepository {
	return &CreateReviewRepository{conn: conn}
}

type CreateReviewRepository struct {
	conn *sqlx.DB
}

func (r *CreateReviewRepository) Add(review *domain.Review) error {
	_, err := r.conn.Exec(
		"INSERT INTO reviews_write (uuid, comment, rating) VALUES (?, ?, ?)",
		review.Uuid(),
		review.Comment(),
		review.Rating(),
	)
	if err != nil {
		return err
	}

	return nil
}
