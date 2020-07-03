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
	"github.com/shinomontaz/chupdate/internal/errorer"
	"github.com/shinomontaz/chupdate/internal/inserter"
	"github.com/shinomontaz/chupdate/internal/parser"

	"github.com/shinomontaz/chupdate/internal/service"
	"github.com/shinomontaz/chupdate/internal/updater"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var env *config.Env
var errors chan error

func init() {
	env = config.NewEnv("./config")
	env.InitDb()
	errors = make(chan error, 1000)
}

func main() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var makeReq func(q, content string, count int)
	makeReq = func(q, content string, count int) {
		fmt.Println("now we send a request", q, content)
	}

	go errorer.Listen(errors)
	parsr := parser.New()
	instr := inserter.New(env.Config.FlushInterval, env.Config.FlushCount, makeReq, errors)
	updtr := updater.New(instr, parsr, env.Db, errors)

	prc := service.New(instr, updtr, parsr, env.Config.CHUrl, env.Db, errors)

	mux := NewMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", env.Config.ListenPort),
		Handler: mux,
	}
	mux.Post("/", prc.Handle)
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
