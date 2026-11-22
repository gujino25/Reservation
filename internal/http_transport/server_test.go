package httptransport_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"reservation/internal/http_transport"
)

func newTestRouter() http.Handler {
	return httptransport.NewRouter(httptransport.Deps{
		Restaurant: httptransport.NewRestaurantHandler(&fakeRestaurantRepo{}),
		Table:      httptransport.NewTableHandler(&fakeTableRepo{}),
		Guest:      httptransport.NewGuestHandler(&fakeGuestRepo{}),
		Booking:    httptransport.NewBookingHandler(&fakeBookingService{}),
	})
}

func TestRouter(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{"create restaurant", http.MethodPost, "/api/v1/restaurants", `{"name":"A","address":"B"}`, http.StatusCreated},
		{"create guest", http.MethodPost, "/api/v1/guests", `{"name":"Ivan"}`, http.StatusCreated},
		{"create booking bad body", http.MethodPost, "/api/v1/bookings", `{`, http.StatusBadRequest},
		{"get restaurant by id", http.MethodGet, "/api/v1/restaurants/3", "", http.StatusNotFound},
		{"get guest by id", http.MethodGet, "/api/v1/guests/3", "", http.StatusNotFound},
		{"unknown path", http.MethodGet, "/api/v1/nope", "", http.StatusNotFound},
		{"method not allowed", http.MethodPut, "/api/v1/restaurants", `{}`, http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			newTestRouter().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d, want %d, body: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}
