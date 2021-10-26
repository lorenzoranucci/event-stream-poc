package application

type GetReviewsQueryHandler struct {
	reviewRepository GetReviewsRepository
}

func NewGetReviewsQueryHandler(reviewRepository GetReviewsRepository) *GetReviewsQueryHandler {
	return &GetReviewsQueryHandler{reviewRepository: reviewRepository}
}

type GetReviewsQuery struct {
	Limit  int32
	Offset int64
}

type Review struct {
	UUID    string
	Comment string
	Rating  int32
}

type GetReviewsRepository interface {
	Find(query GetReviewsQuery) ([]Review, error)
}

func (h *GetReviewsQueryHandler) Execute(query GetReviewsQuery) ([]Review, error) {
	return h.reviewRepository.Find(query)
}
