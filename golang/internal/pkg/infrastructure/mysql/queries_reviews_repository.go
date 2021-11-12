package mysql

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/queries"
)

func NewGetReviewsRepository(db *sqlx.DB) *QueriesReviewsRepository {
	return &QueriesReviewsRepository{db: db}
}

type QueriesReviewsRepository struct {
	db *sqlx.DB
}

type QueryReviewRow struct {
	UUID    string `db:"uuid"`
	Comment string `db:"comment"`
	Rating  int32  `db:"rating"`
}

func (r *QueriesReviewsRepository) Find(query queries.GetReviewsQuery) ([]queries.Review, error) {
	rows, err := r.db.Queryx(
		"SELECT uuid, comment, rating FROM reviews_read LIMIT ? OFFSET ?",
		query.Limit,
		query.Offset,
	)
	if err != nil {
		return nil, err
	}

	results := make([]queries.Review, 0)
	for rows.Next() {
		var r QueryReviewRow
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}

		results = append(
			results,
			queries.Review{
				UUID:    r.UUID,
				Comment: r.Comment,
				Rating:  r.Rating,
			},
		)
	}

	return results, nil
}

func (r *QueriesReviewsRepository) AddOrUpdateByUUID(review *queries.Review) error {
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

func (r *QueriesReviewsRepository) FindByUUID(uuid uuid.UUID) (*queries.Review, error) {
	rows, err := r.db.Queryx(
		"SELECT uuid, comment, rating FROM reviews_read LIMIT WHERE uuid = ?",
		uuid.String(),
	)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var reviewRow *QueryReviewRow
	err = rows.StructScan(reviewRow)
	if err != nil {
		return nil, err
	}

	if reviewRow == nil {
		return nil, nil
	}

	return &queries.Review{
		UUID:    reviewRow.UUID,
		Comment: reviewRow.Comment,
		Rating:  reviewRow.Rating,
	}, nil
}
