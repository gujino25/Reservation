package httptransport

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Deps struct {
	Restaurant *RestaurantHandler
	Table      *TableHandler
	Guest      *GuestHandler
	Booking    *BookingHandler
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(accessLog)
	r.Use(middleware.Recoverer)

	r.Post("/api/v1/restaurants", d.Restaurant.Create)
	r.Get("/api/v1/restaurants/{id}", d.Restaurant.GetByID)
	r.Post("/api/v1/restaurants/{id}/tables", d.Table.Create)
	r.Get("/api/v1/restaurants/{id}/tables", d.Table.ListByRestaurant)

	r.Post("/api/v1/guests", d.Guest.Create)
	r.Get("/api/v1/guests/by-phone", d.Guest.GetByPhone)
	r.Get("/api/v1/guests/{id}", d.Guest.GetByID)

	r.Post("/api/v1/bookings", d.Booking.Create)
	r.Get("/api/v1/bookings/{id}", d.Booking.GetByID)
	r.Post("/api/v1/bookings/{id}/cancel", d.Booking.Cancel)

	return r
}

func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		slog.InfoContext(r.Context(), "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"duration", time.Since(start).String(),
			"request_id", middleware.GetReqID(r.Context()),
		)
	})
}
