package mysql

import (
	"github.com/jmoiron/sqlx"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

func NewCreateReviewRepository(db *sqlx.DB) *CreateReviewRepository {
	return &CreateReviewRepository{db: db}
}

type CreateReviewRepository struct {
	db *sqlx.DB
}

func (r *CreateReviewRepository) Add(review *domain.Review, transaction application.Transaction) error {
	tx, err := getTransaction(transaction, r.db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
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
