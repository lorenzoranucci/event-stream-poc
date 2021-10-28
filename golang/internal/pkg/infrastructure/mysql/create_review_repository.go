package mysql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

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
	tx, shouldCommit, err := getTransaction(transaction, r.db)
	if err != nil {
		return err
	}

	if shouldCommit {
		defer func(tx *sql.Tx) {
			err := tx.Rollback()
			if err != nil {
				logrus.Error(err)
			}
		}(tx)
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

	if shouldCommit {
		return tx.Commit()
	}
	return nil
}
