package httptransport_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"reservation/internal/domain"
	"reservation/internal/http_transport"
)

type fakeRestaurantRepo struct {
	rest        domain.Restaurant
	err         error
	lastCreated domain.Restaurant
}

func (f *fakeRestaurantRepo) Create(_ context.Context, r domain.Restaurant) (domain.Restaurant, error) {
	if f.err != nil {
		return domain.Restaurant{}, f.err
	}
	r.ID = 1
	f.lastCreated = r
	return r, nil
}

func (f *fakeRestaurantRepo) GetByID(_ context.Context, id int64) (domain.Restaurant, error) {
	if f.err != nil {
		return domain.Restaurant{}, f.err
	}
	if f.rest.ID != id {
		return domain.Restaurant{}, domain.ErrNotFound
	}
	return f.rest, nil
}

func TestRestaurantHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repoErr    error
		wantStatus int
	}{
		{name: "success", body: `{"name":"Test Cafe","address":"Main st, 1"}`, wantStatus: http.StatusCreated},
		{name: "invalid json", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "trailing data", body: `{"name":"Test Cafe","address":"Main st, 1"} {}`, wantStatus: http.StatusBadRequest},
		{name: "empty name", body: `{"name":"","address":"Main st, 1"}`, wantStatus: http.StatusBadRequest},
		{name: "repo conflict", body: `{"name":"Test Cafe","address":"Main st, 1"}`, repoErr: domain.ErrAlreadyExists, wantStatus: http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRestaurantRepo{err: tt.repoErr}
			h := httptransport.NewRestaurantHandler(repo)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/restaurants", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			h.Create(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated {
				if repo.lastCreated.Name != "Test Cafe" || repo.lastCreated.Address != "Main st, 1" {
					t.Fatalf("repo got %+v", repo.lastCreated)
				}
				var resp map[string]any
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp["name"] != "Test Cafe" {
					t.Fatalf("response: %v", resp)
				}
			}
		})
	}
}

func TestRestaurantHandler_GetByID(t *testing.T) {
	rest := domain.Restaurant{ID: 7, Name: "Test Cafe", Address: "Main st, 1", CreatedAt: time.Now()}

	t.Run("success", func(t *testing.T) {
		h := httptransport.NewRestaurantHandler(&fakeRestaurantRepo{rest: rest})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants/7", nil)
		req.SetPathValue("id", "7")
		rec := httptest.NewRecorder()
		h.GetByID(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status: got %d, want 200", rec.Code)
		}
		var resp map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if resp["id"] != float64(7) || resp["name"] != "Test Cafe" {
			t.Fatalf("response: %v", resp)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := httptransport.NewRestaurantHandler(&fakeRestaurantRepo{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants/7", nil)
		req.SetPathValue("id", "7")
		rec := httptest.NewRecorder()
		h.GetByID(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status: got %d, want 404", rec.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		h := httptransport.NewRestaurantHandler(&fakeRestaurantRepo{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants/abc", nil)
		req.SetPathValue("id", "abc")
		rec := httptest.NewRecorder()
		h.GetByID(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status: got %d, want 400", rec.Code)
		}
	})
}
