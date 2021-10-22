package domain

import (
	"github.com/google/uuid"
)

func NewReview(comment string, rating int) *Review {
	return &Review{
		uuid:    uuid.New(),
		comment: comment,
	}
}

type Review struct {
	uuid    uuid.UUID
	comment string
	rating  int
}

func (r *Review) Uuid() uuid.UUID {
	return r.uuid
}

func (r *Review) Comment() string {
	return r.comment
}

func (r *Review) Rating() int {
	return r.rating
}
