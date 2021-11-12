package get_reviews

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/queries"
	http2 "github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/utils"
)

type GetReviewsHandler struct {
	qh *queries.GetReviewsQueryHandler
}

func NewGetReviewsHandler(
	getReviewsQueryHandler *queries.GetReviewsQueryHandler,
) *GetReviewsHandler {
	return &GetReviewsHandler{
		qh: getReviewsQueryHandler,
	}
}

type Response []ReviewResponse

type ReviewResponse struct {
	UUID    string `json:"uuid"`
	Comment string `json:"comment"`
	Rating  int32  `json:"rating"`
}

func (h *GetReviewsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Processing request %#v", r)

	if r.Method != "GET" {
		http2.Fail(w, r, fmt.Errorf("only GET method is allowed"), http.StatusMethodNotAllowed)
		return
	}

	limitVal := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitVal)
	if err != nil {
		limit = 10
	}

	offsetVal := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetVal)
	if err != nil {
		offset = 0
	}

	reviews, err := h.qh.Execute(
		queries.GetReviewsQuery{
			Limit:  int32(limit),
			Offset: int64(offset),
		},
	)
	if err != nil {
		logrus.Error(err)
		http2.Fail(w, r, err, http.StatusInternalServerError)
	}

	http2.Send(w, r, createResponse(reviews), http.StatusAccepted)
}

func createResponse(reviews []queries.Review) Response {
	var res = make(Response, len(reviews))
	for i, review := range reviews {
		res[i] = ReviewResponse{
			UUID:    review.UUID,
			Comment: review.Comment,
			Rating:  review.Rating,
		}
	}

	return res
}
