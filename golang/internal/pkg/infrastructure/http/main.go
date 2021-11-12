package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/patch_review"

	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/create_review"
	"github.com/ProntoPro/event-stream-golang/internal/pkg/infrastructure/http/get_reviews"
)

type Server struct {
	m    *http.ServeMux
	port int
}

func NewServer(
	port int,
	createReviewHandler *create_review.CreateReviewHandler,
	getReviewsHandler *get_reviews.GetReviewsHandler,
	incrementReviewRatingHandler *patch_review.PatchReviewHandler,
) *Server {
	r := mux.NewRouter()

	r.Handle("/reviews", createReviewHandler).Methods("POST")
	r.Handle("/reviews", getReviewsHandler).Methods("GET")
	r.Handle("/reviews/{reviewUuid}", incrementReviewRatingHandler).Methods("PATCH")

	srv := &Server{
		m:    http.NewServeMux(),
		port: port,
	}

	srv.m.Handle("/", r)

	return srv
}

func (s *Server) Run() {
	fmt.Println("Server start")

	err := http.ListenAndServe(
		fmt.Sprintf(":%d", s.port),
		s.m,
	)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	fmt.Println("Server terminated")

}
