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

func TestRestaurantRepo_CreateAndGet(t *testing.T) {
	cleanup(t)
	ctx := context.Background()
	repo := postgres.NewRestaurantRepo(testPool)

	created, err := repo.Create(ctx,
		domain.Restaurant{Name: "Test Cafe", Address: "Main st, 1", CreatedAt: time.Now()})
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
	if got.Name != created.Name || got.Address != created.Address {
		t.Fatalf("got %+v, want name/address of %+v", got, created)
	}

	_, err = repo.GetByID(ctx, 999999)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
