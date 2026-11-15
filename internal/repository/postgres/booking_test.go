//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"reservation/internal/domain"
	"reservation/internal/repository/postgres"
)

func TestBookingRepo_CreateAndGet(t *testing.T) {
	cleanup(t)
	ctx := context.Background()
	repo := postgres.NewBookingRepo(testPool)
	tableID, guestID := seed(t)

	start := time.Now().Add(48 * time.Hour).Truncate(time.Microsecond)
	created, err := repo.Create(ctx, newBooking(tableID, guestID, start))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("ID not assigned")
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !got.StartsAt.Equal(start) || !got.EndsAt.Equal(start.Add(2*time.Hour)) {
		t.Fatalf("times: got %v–%v, want %v–%v",
			got.StartsAt, got.EndsAt, start, start.Add(2*time.Hour))
	}
	if got.Status != domain.ReservationStatusActive || got.GuestsAmount != 2 {
		t.Fatalf("got %+v", got)
	}

	_, err = repo.GetByID(ctx, 999999)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestBookingRepo_Overlap(t *testing.T) {
	cleanup(t)
	ctx := context.Background()
	repo := postgres.NewBookingRepo(testPool)
	tableID, guestID := seed(t)

	base := time.Now().Add(48 * time.Hour).Truncate(time.Microsecond)
	if _, err := repo.Create(ctx, newBooking(tableID, guestID, base)); err != nil {
		t.Fatalf("seed booking: %v", err)
	}

	conflicts := []struct {
		name  string
		start time.Time
	}{
		{"exact match", base},
		{"inside", base.Add(30 * time.Minute)},
		{"tail overlap", base.Add(-time.Hour)},
		{"covers whole", base.Add(-30 * time.Minute)},
	}
	for _, tt := range conflicts {
		t.Run("conflict/"+tt.name, func(t *testing.T) {
			_, err := repo.Create(ctx, newBooking(tableID, guestID, tt.start))
			if !errors.Is(err, domain.ErrTableBusy) {
				t.Fatalf("want ErrTableBusy, got %v", err)
			}
		})
	}

	adjacent := []struct {
		name  string
		start time.Time
	}{
		{"right", base.Add(2 * time.Hour)},
		{"left", base.Add(-2 * time.Hour)},
	}
	for _, tt := range adjacent {
		t.Run("adjacent/"+tt.name, func(t *testing.T) {
			if _, err := repo.Create(ctx, newBooking(tableID, guestID, tt.start)); err != nil {
				t.Fatalf("adjacent booking must pass, got %v", err)
			}
		})
	}
}

func TestBookingRepo_CancelledDoesNotBlock(t *testing.T) {
	cleanup(t)
	ctx := context.Background()
	repo := postgres.NewBookingRepo(testPool)
	tableID, guestID := seed(t)

	base := time.Now().Add(48 * time.Hour).Truncate(time.Microsecond)

	b, err := repo.Create(ctx, newBooking(tableID, guestID, base))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	b.Status = domain.ReservationStatusCancelled
	if err := repo.Update(ctx, b); err != nil {
		t.Fatalf("cancel: %v", err)
	}

	if _, err := repo.Create(ctx, newBooking(tableID, guestID, base)); err != nil {
		t.Fatalf("cancelled booking must not block the slot: %v", err)
	}

	// и наоборот: активную бронь нельзя передвинуть на чужое время
	b2, err := repo.Create(ctx, newBooking(tableID, guestID, base.Add(4*time.Hour)))
	if err != nil {
		t.Fatalf("create b2: %v", err)
	}
	b2.StartsAt = base.Add(time.Hour)
	b2.EndsAt = base.Add(3 * time.Hour)
	if err := repo.Update(ctx, b2); !errors.Is(err, domain.ErrTableBusy) {
		t.Fatalf("moving into busy slot: want ErrTableBusy, got %v", err)
	}
}

func TestBookingRepo_UpdateNotFound(t *testing.T) {
	cleanup(t)
	err := postgres.NewBookingRepo(testPool).Update(context.Background(), domain.Reservation{ID: 12345})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
