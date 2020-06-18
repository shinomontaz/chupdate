package service

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
)

type Service struct {
}

// func New(flush_interval, flush_count int) *Service {
// 	return nil
// }

func New(c Collector, db *sqlx.DB) *Service {
	return nil
}

func (s *Service) Process(w http.ResponseWriter, r *http.Request) {
	// parse query
	// return result

	w.WriteHeader(http.StatusOK)
}

func (s *Service) Shutdown(ctx context.Context) error {
	log.Printf("service shutting down\n")
	return nil
}
