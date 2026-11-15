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

func TestTableRepo_ListAndUnique(t *testing.T) {
	cleanup(t)
	ctx := context.Background()
	repo := postgres.NewTableRepo(testPool)
	restRepo := postgres.NewRestaurantRepo(testPool)

	r1, err := restRepo.Create(ctx, domain.Restaurant{Name: "R1", Address: "a", CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("create restaurant 1: %v", err)
	}
	r2, err := restRepo.Create(ctx, domain.Restaurant{Name: "R2", Address: "b", CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("create restaurant 2: %v", err)
	}

	t1, err := repo.Create(ctx, mustTable(t, r1.ID, 1, 4))
	if err != nil {
		t.Fatalf("create table 1: %v", err)
	}
	t2, err := repo.Create(ctx, mustTable(t, r1.ID, 2, 2))
	if err != nil {
		t.Fatalf("create table 2: %v", err)
	}
	if _, err := repo.Create(ctx, mustTable(t, r2.ID, 1, 6)); err != nil {
		t.Fatalf("create table for r2: %v", err)
	}

	tables, err := repo.ListByRestaurant(ctx, r1.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(tables) != 2 {
		t.Fatalf("want 2 tables for r1, got %d: %+v", len(tables), tables)
	}
	if tables[0].ID != t1.ID || tables[1].ID != t2.ID {
		t.Fatalf("want order by number [%d %d], got [%d %d]", t1.ID, t2.ID, tables[0].ID, tables[1].ID)
	}

	_, err = repo.Create(ctx, mustTable(t, r1.ID, 1, 8))
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("duplicate (restaurant_id, number): want ErrAlreadyExists, got %v", err)
	}
}
