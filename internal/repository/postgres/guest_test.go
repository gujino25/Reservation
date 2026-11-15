//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"reservation/internal/domain"
	"reservation/internal/repository/postgres"
)

func TestGuestRepo_PhoneUnique(t *testing.T) {
	cleanup(t)
	ctx := context.Background()
	repo := postgres.NewGuestRepo(testPool)

	g, err := repo.Create(ctx, domain.NewGuest("Ivan", "+79001234567"))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	byID, err := repo.GetByID(ctx, g.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if byID.Name != "Ivan" || byID.Phone != "+79001234567" {
		t.Fatalf("got %+v", byID)
	}

	byPhone, err := repo.GetByPhone(ctx, "+79001234567")
	if err != nil {
		t.Fatalf("get by phone: %v", err)
	}
	if byPhone.ID != g.ID {
		t.Fatalf("get by phone: got id %d, want %d", byPhone.ID, g.ID)
	}

	if _, err := repo.Create(ctx, domain.NewGuest("Petr", "+79001234567")); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("duplicate phone: want ErrAlreadyExists, got %v", err)
	}

	// пустой телефон хранится как NULL — дубликатов быть не должно
	if _, err := repo.Create(ctx, domain.NewGuest("NoPhone1", "")); err != nil {
		t.Fatalf("empty phone 1: %v", err)
	}
	if _, err := repo.Create(ctx, domain.NewGuest("NoPhone2", "")); err != nil {
		t.Fatalf("empty phone 2: %v", err)
	}

	_, err = repo.GetByPhone(ctx, "+79999999999")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("missing phone: want ErrNotFound, got %v", err)
	}
}
