package create_review

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"

	http2 "github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/utils"
)

type CreateReviewHandler struct {
	ch *commands.CreateReviewCommandHandler
}

func NewCreateReviewHandler(
	createReviewCommandHandler *commands.CreateReviewCommandHandler,
) *CreateReviewHandler {
	return &CreateReviewHandler{
		ch: createReviewCommandHandler,
	}
}

type Request struct {
	Comment string `json:"comment"`
	Rating  int32  `json:"rating"`
}

func (h *CreateReviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Processing request %#v", r)

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
		commands.CreateReviewCommand{
			Comment: request.Comment,
			Rating:  request.Rating,
		},
	)
	if err != nil {
		logrus.Error(err)
		http2.Fail(w, r, err, http.StatusInternalServerError)
	}

	http2.Send(w, r, struct{}{}, http.StatusAccepted)
}
