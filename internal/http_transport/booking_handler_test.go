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

type fakeBookingService struct {
	booking   domain.Reservation
	err       error
	created   domain.Reservation
	cancelled domain.Reservation
}

func (f *fakeBookingService) Create(_ context.Context, tableID, guestID int64, startsAt time.Time, guests int) (domain.Reservation, error) {
	if f.err != nil {
		return domain.Reservation{}, f.err
	}
	b := f.booking
	b.ID = 1
	b.TableID, b.GuestID, b.StartsAt, b.GuestsAmount = tableID, guestID, startsAt, guests
	f.created = b
	return b, nil
}

func (f *fakeBookingService) GetByID(_ context.Context, id int64) (domain.Reservation, error) {
	if f.err != nil {
		return domain.Reservation{}, f.err
	}
	if f.booking.ID != id {
		return domain.Reservation{}, domain.ErrNotFound
	}
	return f.booking, nil
}

func (f *fakeBookingService) Cancel(_ context.Context, id int64) (domain.Reservation, error) {
	if f.err != nil {
		return domain.Reservation{}, f.err
	}
	if f.booking.ID != id {
		return domain.Reservation{}, domain.ErrNotFound
	}
	b := f.booking
	b.Status = domain.ReservationStatusCancelled
	f.cancelled = b
	return b, nil
}

const validBookingBody = `{"table_id":1,"guest_id":2,"starts_at":"2026-12-01T18:00:00Z","guests_count":2}`

func TestBookingHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcErr     error
		wantStatus int
	}{
		{name: "success", body: validBookingBody, wantStatus: http.StatusCreated},
		{name: "invalid json", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "zero table id", body: `{"table_id":0,"guest_id":2,"starts_at":"2026-12-01T18:00:00Z","guests_count":2}`, wantStatus: http.StatusBadRequest},
		{name: "table busy", body: validBookingBody, svcErr: domain.ErrTableBusy, wantStatus: http.StatusConflict},
		{name: "capacity exceeded", body: validBookingBody, svcErr: domain.ErrCapacityExceeded, wantStatus: http.StatusConflict},
		{name: "table not found", body: validBookingBody, svcErr: domain.ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "start in the past", body: validBookingBody, svcErr: domain.ErrReservationTime, wantStatus: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeBookingService{err: tt.svcErr}
			h := httptransport.NewBookingHandler(svc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			h.Create(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated {
				if svc.created.TableID != 1 || svc.created.GuestsAmount != 2 {
					t.Fatalf("svc got %+v", svc.created)
				}
			}
		})
	}
}

func TestBookingHandler_GetByID(t *testing.T) {
	booking := domain.Reservation{ID: 9, TableID: 1, GuestID: 2, GuestsAmount: 2, Status: domain.ReservationStatusActive}

	t.Run("success", func(t *testing.T) {
		h := httptransport.NewBookingHandler(&fakeBookingService{booking: booking})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings/9", nil)
		req.SetPathValue("id", "9")
		rec := httptest.NewRecorder()
		h.GetByID(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status: got %d, want 200", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := httptransport.NewBookingHandler(&fakeBookingService{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings/9", nil)
		req.SetPathValue("id", "9")
		rec := httptest.NewRecorder()
		h.GetByID(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status: got %d, want 404", rec.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		h := httptransport.NewBookingHandler(&fakeBookingService{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings/abc", nil)
		req.SetPathValue("id", "abc")
		rec := httptest.NewRecorder()
		h.GetByID(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status: got %d, want 400", rec.Code)
		}
	})
}

func TestBookingHandler_Cancel(t *testing.T) {
	booking := domain.Reservation{ID: 9, TableID: 1, GuestID: 2, GuestsAmount: 2, Status: domain.ReservationStatusActive}

	t.Run("success", func(t *testing.T) {
		svc := &fakeBookingService{booking: booking}
		h := httptransport.NewBookingHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/9/cancel", nil)
		req.SetPathValue("id", "9")
		rec := httptest.NewRecorder()
		h.Cancel(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status: got %d, want 200", rec.Code)
		}
		if svc.cancelled.Status != domain.ReservationStatusCancelled {
			t.Fatalf("cancelled status: got %q", svc.cancelled.Status)
		}
	})

	t.Run("already cancelled", func(t *testing.T) {
		h := httptransport.NewBookingHandler(&fakeBookingService{err: domain.ErrReservationAlreadyCancelled})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/9/cancel", nil)
		req.SetPathValue("id", "9")
		rec := httptest.NewRecorder()
		h.Cancel(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf("status: got %d, want 409", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := httptransport.NewBookingHandler(&fakeBookingService{})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/9/cancel", nil)
		req.SetPathValue("id", "9")
		rec := httptest.NewRecorder()
		h.Cancel(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status: got %d, want 404", rec.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		h := httptransport.NewBookingHandler(&fakeBookingService{})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/abc/cancel", nil)
		req.SetPathValue("id", "abc")
		rec := httptest.NewRecorder()
		h.Cancel(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status: got %d, want 400", rec.Code)
		}
	})
}
