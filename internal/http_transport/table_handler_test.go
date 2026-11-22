package httptransport_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"reservation/internal/domain"
	"reservation/internal/http_transport"
)

type fakeTableRepo struct {
	tables     []domain.Table
	err        error
	created    domain.Table
	listRestID int64
}

func (f *fakeTableRepo) Create(_ context.Context, tbl domain.Table) (domain.Table, error) {
	if f.err != nil {
		return domain.Table{}, f.err
	}
	tbl.ID = 1
	f.created = tbl
	return tbl, nil
}

func (f *fakeTableRepo) ListByRestaurant(_ context.Context, restaurantID int64) ([]domain.Table, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.listRestID = restaurantID
	return f.tables, nil
}

func TestTableHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repoErr    error
		wantStatus int
	}{
		{name: "success", body: `{"number":1,"capacity":4}`, wantStatus: http.StatusCreated},
		{name: "invalid json", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "zero capacity", body: `{"number":1,"capacity":0}`, wantStatus: http.StatusBadRequest},
		{name: "zero number", body: `{"number":0,"capacity":4}`, wantStatus: http.StatusBadRequest},
		{name: "restaurant not found", body: `{"number":1,"capacity":4}`, repoErr: domain.ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "duplicate table number", body: `{"number":1,"capacity":4}`, repoErr: domain.ErrAlreadyExists, wantStatus: http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeTableRepo{err: tt.repoErr}
			h := httptransport.NewTableHandler(repo)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/restaurants/5/tables", strings.NewReader(tt.body))
			req.SetPathValue("id", "5")
			rec := httptest.NewRecorder()
			h.Create(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated {
				if repo.created.RestaurantID != 5 || repo.created.Capacity != 4 {
					t.Fatalf("repo got %+v", repo.created)
				}
			}
		})
	}
}

func TestTableHandler_ListByRestaurant(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &fakeTableRepo{tables: []domain.Table{
			{ID: 1, RestaurantID: 5, NumberOfTable: 1, Capacity: 4},
			{ID: 2, RestaurantID: 5, NumberOfTable: 2, Capacity: 2},
		}}
		h := httptransport.NewTableHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants/5/tables", nil)
		req.SetPathValue("id", "5")
		rec := httptest.NewRecorder()
		h.ListByRestaurant(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status: got %d, want 200", rec.Code)
		}
		if repo.listRestID != 5 {
			t.Fatalf("repo got restaurantID %d, want 5", repo.listRestID)
		}
		var resp []map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if len(resp) != 2 {
			t.Fatalf("want 2 tables, got %d", len(resp))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		h := httptransport.NewTableHandler(&fakeTableRepo{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants/5/tables", nil)
		req.SetPathValue("id", "5")
		rec := httptest.NewRecorder()
		h.ListByRestaurant(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status: got %d, want 200", rec.Code)
		}
		if strings.TrimSpace(rec.Body.String()) != "[]" {
			t.Fatalf("want empty JSON array, got %s", rec.Body.String())
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		h := httptransport.NewTableHandler(&fakeTableRepo{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/restaurants/abc/tables", nil)
		req.SetPathValue("id", "abc")
		rec := httptest.NewRecorder()
		h.ListByRestaurant(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status: got %d, want 400", rec.Code)
		}
	})
}
