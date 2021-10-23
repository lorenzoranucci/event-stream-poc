package domain

import (
	"github.com/google/uuid"
)

func NewReview(comment string, rating int32) *Review {
	return &Review{
		uuid:    uuid.New(),
		comment: comment,
		rating: rating,
	}
}

type Review struct {
	uuid    uuid.UUID
	comment string
	rating  int32
}

func (r *Review) Uuid() uuid.UUID {
	return r.uuid
}

func (r *Review) Comment() string {
	return r.comment
}

func (r *Review) Rating() int32 {
	return r.rating
}
