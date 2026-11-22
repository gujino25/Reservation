package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"reservation/internal/domain"
	"reservation/internal/service"
)

type fakeTableRepo struct {
	table domain.Table
	err   error
}

func (f fakeTableRepo) GetByID(_ context.Context, id int64) (domain.Table, error) {
	if f.err != nil {
		return domain.Table{}, f.err
	}
	if f.table.ID != id {
		return domain.Table{}, domain.ErrNotFound
	}
	return f.table, nil
}

type fakeGuestRepo struct {
	guest domain.Guest
	err   error
}

func (f fakeGuestRepo) GetByID(_ context.Context, id int64) (domain.Guest, error) {
	if f.err != nil {
		return domain.Guest{}, f.err
	}
	if f.guest.ID != id {
		return domain.Guest{}, domain.ErrNotFound
	}
	return f.guest, nil
}

type fakeBookingRepo struct {
	createErr error

	byID      map[int64]domain.Reservation
	updated   domain.Reservation
	updateErr error
}

func newFakeBookingRepo() *fakeBookingRepo {
	return &fakeBookingRepo{byID: make(map[int64]domain.Reservation)}
}

func (f *fakeBookingRepo) Create(_ context.Context, b domain.Reservation) (domain.Reservation, error) {
	if f.createErr != nil {
		return domain.Reservation{}, f.createErr
	}
	b.ID = int64(len(f.byID)) + 1
	f.byID[b.ID] = b
	return b, nil
}

func (f *fakeBookingRepo) GetByID(_ context.Context, id int64) (domain.Reservation, error) {
	b, ok := f.byID[id]
	if !ok {
		return domain.Reservation{}, domain.ErrNotFound
	}
	return b, nil
}

func (f *fakeBookingRepo) Update(_ context.Context, b domain.Reservation) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	if _, ok := f.byID[b.ID]; !ok {
		return domain.ErrNotFound
	}
	f.updated = b
	f.byID[b.ID] = b
	return nil
}

func (f *fakeBookingRepo) Cancel(_ context.Context, id int64) error {
	b, ok := f.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	if b.Status != domain.ReservationStatusActive {
		return domain.ErrReservationAlreadyCancelled
	}
	b.Status = domain.ReservationStatusCancelled
	f.updated = b
	f.byID[id] = b
	return nil
}

func TestBookingService_Create(t *testing.T) {
	future := time.Now().Add(48 * time.Hour)
	table := domain.Table{ID: 1, RestaurantID: 1, NumberOfTable: 1, Capacity: 4}
	guest := domain.Guest{ID: 2, Name: "Ivan", Phone: "+79001234567"}

	tests := []struct {
		name      string
		table     domain.Table
		tableErr  error
		guestErr  error
		createErr error
		start     time.Time
		guests    int
		wantErr   error
	}{
		{name: "success", table: table, start: future, guests: 2},
		{name: "table not found", tableErr: domain.ErrNotFound, start: future, guests: 2, wantErr: domain.ErrNotFound},
		{name: "guest not found", table: table, guestErr: domain.ErrNotFound, start: future, guests: 2, wantErr: domain.ErrNotFound},
		{name: "capacity exceeded", table: domain.Table{ID: 1, Capacity: 2}, start: future, guests: 5, wantErr: domain.ErrCapacityExceeded},
		{name: "table busy", table: table, createErr: domain.ErrTableBusy, start: future, guests: 2, wantErr: domain.ErrTableBusy},
		{name: "start in the past", table: table, start: time.Now().Add(-time.Hour), guests: 2, wantErr: domain.ErrReservationTime},
		{name: "zero guests", table: table, start: future, guests: 0, wantErr: domain.ErrGuestAmountInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookings := newFakeBookingRepo()
			bookings.createErr = tt.createErr
			svc := service.NewBookingService(
				fakeTableRepo{table: tt.table, err: tt.tableErr},
				fakeGuestRepo{guest: guest, err: tt.guestErr},
				bookings,
			)

			got, err := svc.Create(context.Background(), tt.table.ID, guest.ID, tt.start, tt.guests)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("want %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID == 0 {
				t.Fatal("ID not assigned")
			}
			if got.Status != domain.ReservationStatusActive {
				t.Fatalf("status: got %q, want %q", got.Status, domain.ReservationStatusActive)
			}
			if !got.EndsAt.Equal(got.StartsAt.Add(2 * time.Hour)) {
				t.Fatalf("ends: got %v, want starts+2h", got.EndsAt)
			}
		})
	}
}

func TestBookingService_Cancel(t *testing.T) {
	future := time.Now().Add(48 * time.Hour)
	active := domain.Reservation{
		ID: 1, TableID: 1, GuestID: 2, GuestsAmount: 2,
		CreatedAt: time.Now(),
		StartsAt:  future, EndsAt: future.Add(2 * time.Hour),
		Status: domain.ReservationStatusActive,
	}
	cancelled := active
	cancelled.Status = domain.ReservationStatusCancelled
	started := active
	started.ID = 3
	started.StartsAt = time.Now().Add(-time.Hour)
	started.EndsAt = started.StartsAt.Add(2 * time.Hour)

	tests := []struct {
		name    string
		booking domain.Reservation
		seed    bool
		wantErr error
	}{
		{name: "success", booking: active, seed: true},
		{name: "not found", booking: domain.Reservation{ID: 99}, wantErr: domain.ErrNotFound},
		{name: "already cancelled", booking: cancelled, seed: true, wantErr: domain.ErrReservationAlreadyCancelled},
		{name: "already started", booking: started, seed: true, wantErr: domain.ErrReservationAlreadyStarted},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookings := newFakeBookingRepo()
			if tt.seed {
				bookings.byID[tt.booking.ID] = tt.booking
			}
			svc := service.NewBookingService(fakeTableRepo{}, fakeGuestRepo{}, bookings)

			got, err := svc.Cancel(context.Background(), tt.booking.ID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("want %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Status != domain.ReservationStatusCancelled {
				t.Fatalf("returned status: got %q, want %q", got.Status, domain.ReservationStatusCancelled)
			}
			if bookings.updated.Status != domain.ReservationStatusCancelled {
				t.Fatalf("saved status: got %q, want %q", bookings.updated.Status, domain.ReservationStatusCancelled)
			}
		})
	}
}
