package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"reservation/internal/config"
	httptransport "reservation/internal/http_transport"
	"reservation/internal/repository/postgres"
	"reservation/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := postgres.Migrate(cfg.DatabaseURL); err != nil {
		return err
	}
	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	tableRepo := postgres.NewTableRepo(pool)
	guestRepo := postgres.NewGuestRepo(pool)
	bookingSvc := service.NewBookingService(tableRepo, guestRepo, postgres.NewBookingRepo(pool))

	router := httptransport.NewRouter(httptransport.Deps{
		Restaurant: httptransport.NewRestaurantHandler(postgres.NewRestaurantRepo(pool)),
		Table:      httptransport.NewTableHandler(tableRepo),
		Guest:      httptransport.NewGuestHandler(guestRepo),
		Booking:    httptransport.NewBookingHandler(bookingSvc),
	})
	srv := httptransport.NewServer(cfg.HTTPAddr, router)

	errCh := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", cfg.HTTPAddr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
