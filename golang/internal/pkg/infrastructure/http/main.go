package http

import (
	"fmt"
	"net/http"

	"github.com/ProntoPro/golang-kafka/internal/pkg/infrastructure/http/create_review"
)

type Server struct {
	m           *http.ServeMux
	port        int
}

func NewServer(
	port int,
	createReviewHandler *create_review.CreateReviewHandler,
) *Server {
	srv := &Server{
		m: http.NewServeMux(),
		port: port,
	}

	srv.m.Handle("/reviews", createReviewHandler)

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
