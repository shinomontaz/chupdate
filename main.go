package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shinomontaz/chupdate/config"
	"github.com/shinomontaz/chupdate/internal/collector"
	"github.com/shinomontaz/chupdate/internal/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	_ "github.com/ClickHouse/clickhouse-go"
)

var env *config.Env

func init() {
	env = config.NewEnv("./config")
	env.InitDb()
}

func main() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	collr := collector.New(env.Config.FlushInterval, env.Config.FlushCount)
	prc := service.New(collr, env.Db)

	mux := NewMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", env.Config.ListenPort),
		Handler: mux,
	}
	mux.Post("/", prc.Process)
	mux.Get("/metrics", promhttp.Handler().ServeHTTP)
	//	mux.Get("/alive", alive handler)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Debug("server started on port: ", env.Config.ListenPort)

	<-signals

	log.Debug("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		prc.Shutdown(ctx)
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error %+v\n", err)

		os.Exit(1)
	}
}

func NewMux() *chi.Mux {
	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(
		middleware.Compress(5, "gzip"),
		middleware.RedirectSlashes,
		middleware.Recoverer,
		cors.Handler,
	)
	//	router.Use(middleware.Logger)

	return router
}
