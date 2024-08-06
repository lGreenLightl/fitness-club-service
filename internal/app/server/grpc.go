package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RunGRPCServer(createHandler func(router chi.Router) http.Handler) {
	// TODO: init grpc server
}