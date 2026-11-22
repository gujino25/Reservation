package httptransport_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"reservation/internal/domain"
	"reservation/internal/http_transport"
)

type fakeGuestRepo struct {
	guest   domain.Guest
	err     error
	created domain.Guest
}

func (f *fakeGuestRepo) Create(_ context.Context, g domain.Guest) (domain.Guest, error) {
	if f.err != nil {
		return domain.Guest{}, f.err
	}
	g.ID = 1
	f.created = g
	return g, nil
}

func (f *fakeGuestRepo) GetByID(_ context.Context, id int64) (domain.Guest, error) {
	if f.err != nil {
		return domain.Guest{}, f.err
	}
	if f.guest.ID != id {
		return domain.Guest{}, domain.ErrNotFound
	}
	return f.guest, nil
}

func (f *fakeGuestRepo) GetByPhone(_ context.Context, phone string) (domain.Guest, error) {
	if f.err != nil {
		return domain.Guest{}, f.err
	}
	if f.guest.Phone != phone {
		return domain.Guest{}, domain.ErrNotFound
	}
	return f.guest, nil
}

func TestGuestHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repoErr    error
		wantStatus int
	}{
		{name: "success", body: `{"name":"Ivan","phone":"+79001234567"}`, wantStatus: http.StatusCreated},
		{name: "invalid json", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "empty name", body: `{"name":""}`, wantStatus: http.StatusBadRequest},
		{name: "duplicate phone", body: `{"name":"Ivan","phone":"+79001234567"}`, repoErr: domain.ErrAlreadyExists, wantStatus: http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeGuestRepo{err: tt.repoErr}
			h := httptransport.NewGuestHandler(repo)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/guests", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			h.Create(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated && repo.created.Name != "Ivan" {
				t.Fatalf("repo got %+v", repo.created)
			}
		})
	}
}

func TestGuestHandler_GetByID(t *testing.T) {
	guest := domain.Guest{ID: 3, Name: "Ivan", Phone: "+79001234567", CreatedAt: time.Now()}

	t.Run("success", func(t *testing.T) {
		h := httptransport.NewGuestHandler(&fakeGuestRepo{guest: guest})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/guests/3", nil)
		req.SetPathValue("id", "3")
		rec := httptest.NewRecorder()
		h.GetByID(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status: got %d, want 200", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := httptransport.NewGuestHandler(&fakeGuestRepo{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/guests/3", nil)
		req.SetPathValue("id", "3")
		rec := httptest.NewRecorder()
		h.GetByID(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status: got %d, want 404", rec.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		h := httptransport.NewGuestHandler(&fakeGuestRepo{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/guests/abc", nil)
		req.SetPathValue("id", "abc")
		rec := httptest.NewRecorder()
		h.GetByID(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status: got %d, want 400", rec.Code)
		}
	})
}

func TestGuestHandler_GetByPhone(t *testing.T) {
	guest := domain.Guest{ID: 3, Name: "Ivan", Phone: "+79001234567", CreatedAt: time.Now()}

	t.Run("success", func(t *testing.T) {
		h := httptransport.NewGuestHandler(&fakeGuestRepo{guest: guest})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/guests/by-phone?phone=%2B79001234567", nil)
		rec := httptest.NewRecorder()
		h.GetByPhone(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status: got %d, want 200", rec.Code)
		}
	})

	t.Run("missing phone param", func(t *testing.T) {
		h := httptransport.NewGuestHandler(&fakeGuestRepo{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/guests/by-phone", nil)
		rec := httptest.NewRecorder()
		h.GetByPhone(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status: got %d, want 400", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := httptransport.NewGuestHandler(&fakeGuestRepo{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/guests/by-phone?phone=%2B70000000000", nil)
		rec := httptest.NewRecorder()
		h.GetByPhone(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status: got %d, want 404", rec.Code)
		}
	})
}
