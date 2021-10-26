package mysql

import (
	"github.com/jmoiron/sqlx"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application"
)

func NewGetReviewsRepository(db *sqlx.DB) *GetReviewsRepository {
	return &GetReviewsRepository{db: db}
}

type GetReviewsRepository struct {
	db *sqlx.DB
}

type ReviewRow struct {
	UUID    string `db:"uuid"`
	Comment string `db:"comment"`
	Rating  int32  `db:"rating"`
}

func (r *GetReviewsRepository) Find(query application.GetReviewsQuery) ([]application.Review, error) {
	rows, err := r.db.Queryx(
		"SELECT uuid, comment, rating FROM reviews_read LIMIT ? OFFSET ?",
		query.Limit,
		query.Offset,
	)
	if err != nil {
		return nil, err
	}

	results := make([]application.Review, 0)
	for rows.Next() {
		var r ReviewRow
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}

		results = append(
			results,
			application.Review{
				UUID:    r.UUID,
				Comment: r.Comment,
				Rating:  r.Rating,
			},
		)
	}

	return results, nil
}

func (r *GetReviewsRepository) AddOrUpdateByUUID(review *application.Review) error {
	_, err := r.db.Exec(
		"INSERT INTO reviews_read (uuid, comment, rating) VALUES (?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE comment=?, rating=?",
		review.UUID,
		review.Comment,
		review.Rating,
		review.Comment,
		review.Rating,
	)
	if err != nil {
		return err
	}

	return nil
}
