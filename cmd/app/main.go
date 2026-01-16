package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"github.com/MichaelShi-san/subscription-test/internal/config"
	"github.com/MichaelShi-san/subscription-test/internal/handler"
	"github.com/MichaelShi-san/subscription-test/internal/logger"
	"github.com/MichaelShi-san/subscription-test/internal/repository"
	"github.com/MichaelShi-san/subscription-test/internal/service"
)

func main() {
	cfg := config.MustLoad()

	logg := logger.New()

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewSubscriptionRepository(db)
	svc := service.NewSubscriptionService(repo)
	h := handler.NewSubscriptionHandler(svc, logg)

	r := chi.NewRouter()

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/total", h.TotalCost)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	go func() {
		logg.Info("http server started", "port", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logg.Info("shutting down server")
	_ = srv.Shutdown(ctx)
}
