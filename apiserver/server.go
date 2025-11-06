package apiserver

import (
	"context"
	"log"
	"net/http"

	"github.com/dankski/learn-asyncapi/config"
)

type ApiServer struct {
	Config *config.Config
}

func New(config *config.Config) *ApiServer {
	return &ApiServer{Config: config}
}

func (s *ApiServer) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (s *ApiServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", s.ping)

	server := &http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	log.Printf("starting server on 5000...\n")

	return server.ListenAndServe()
}
