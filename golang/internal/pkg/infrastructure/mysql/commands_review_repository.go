package mysql

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/domain"
)

func NewReviewRepository(db *sqlx.DB) *CommandsReviewRepository {
	return &CommandsReviewRepository{db: db}
}

type CommandsReviewRepository struct {
	db *sqlx.DB
}

type CommandReviewRow struct {
	ID      int64  `db:"id"`
	UUID    string `db:"uuid"`
	Comment string `db:"comment"`
	Rating  int32  `db:"rating"`
}

func (r *CommandsReviewRepository) FindByUUID(reviewUUID uuid.UUID) (*domain.Review, error) {
	rows, err := r.db.Queryx(
		"SELECT id, uuid, comment, rating FROM reviews_read WHERE uuid = ?",
		reviewUUID.String(),
	)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	if !next {
		return nil, nil
	}

	var reviewRow *CommandReviewRow
	err = rows.StructScan(reviewRow)
	if err != nil {
		return nil, err
	}

	if reviewRow == nil {
		return nil, nil
	}

	return domain.CreateFromRepository(reviewUUID, reviewRow.Comment, reviewRow.Rating), nil
}

func (r *CommandsReviewRepository) Save(review *domain.Review, transaction commands.Transaction) error {
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
