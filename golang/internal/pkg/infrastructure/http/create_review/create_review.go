package create_review

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ProntoPro/golang-kafka/internal/pkg/application"
	http2 "github.com/ProntoPro/golang-kafka/internal/pkg/infrastructure/http/utils"
)

type CreateReviewHandler struct {
	ch *application.CreateReviewCommandHandler
}

func NewCreateReviewHandler(
	createReviewCommandHandler *application.CreateReviewCommandHandler,
) *CreateReviewHandler {
	return &CreateReviewHandler{
		ch: createReviewCommandHandler,
	}
}

type Request struct {
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
}

func (h *CreateReviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http2.Fail(w, r, fmt.Errorf("only POST method is allowed"), http.StatusMethodNotAllowed)
		return
	}

	request := &Request{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		http2.Fail(w, r, err, http.StatusUnprocessableEntity)
	}

	err = h.ch.Execute(
		application.CreateReviewCommand{
			Comment: request.Comment,
			Rating:  request.Rating,
		},
	)
	if err != nil {
		http2.Fail(w, r, err, http.StatusInternalServerError)
	}

	http2.Send(w, r, struct{}{}, http.StatusAccepted)
}
