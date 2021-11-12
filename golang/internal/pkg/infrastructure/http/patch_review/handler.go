package patch_review

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/application/commands"

	http2 "github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/utils"
)

type PatchReviewHandler struct {
	ch *commands.IncrementReviewRatingCommandHandler
}

func NewPatchReviewHandler(
	patchReviewHandler *commands.IncrementReviewRatingCommandHandler,
) *PatchReviewHandler {
	return &PatchReviewHandler{
		ch: patchReviewHandler,
	}
}

// ServeHTTP is not RESTfulish and well-designed, but it's acceptable for the POC goal
func (h *PatchReviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Processing request %#v", r)

	if r.Method != http.MethodPatch {
		http2.Fail(w, r, fmt.Errorf("only PATCH method is allowed"), http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	reviewUUIDString, found := vars["reviewUuid"]
	if !found {
		http2.Fail(w, r, fmt.Errorf("missing review uuid"), http.StatusBadRequest)
		return
	}

	reviewUUID, err := uuid.Parse(reviewUUIDString)
	if err != nil {
		http2.Fail(w, r, fmt.Errorf("wrong review uuid"), http.StatusBadRequest)
		return
	}

	err = h.ch.Execute(
		commands.IncrementReviewRatingCommand{
			ReviewUUID: reviewUUID,
		},
	)
	if err != nil {
		logrus.Error(err)
		http2.Fail(w, r, err, http.StatusInternalServerError)
	}

	http2.Send(w, r, struct{}{}, http.StatusAccepted)
}
