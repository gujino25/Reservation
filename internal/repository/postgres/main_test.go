//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"reservation/internal/domain"
	"reservation/internal/repository/postgres"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	ctx := context.Background()

	ctr, err := tcpostgres.Run(ctx, "postgres:18",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres container: %v\n", err)
		return 1
	}
	defer func() {
		if err := ctr.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "terminate container: %v\n", err)
		}
	}()

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection string: %v\n", err)
		return 1
	}
	if err := postgres.Migrate(dsn); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		return 1
	}
	pool, err := postgres.NewPool(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pool: %v\n", err)
		return 1
	}
	defer pool.Close()

	testPool = pool
	return m.Run()
}

func cleanup(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		`TRUNCATE bookings, guests, tables, restaurants CASCADE`)
	if err != nil {
		t.Fatalf("cleanup: %v", err)
	}
}

func mustTable(t *testing.T, restaurantID int64, number, capacity int) domain.Table {
	t.Helper()
	tbl, err := domain.NewTable(restaurantID, number, capacity)
	if err != nil {
		t.Fatalf("NewTable: %v", err)
	}
	return tbl
}

// создаёт ресторан, стол на 4 места и гостя
func seed(t *testing.T) (tableID, guestID int64) {
	t.Helper()
	ctx := context.Background()

	rest, err := postgres.NewRestaurantRepo(testPool).Create(ctx,
		domain.Restaurant{Name: "Test Cafe", Address: "Main st, 1", CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("seed restaurant: %v", err)
	}
	tbl, err := postgres.NewTableRepo(testPool).Create(ctx, mustTable(t, rest.ID, 1, 4))
	if err != nil {
		t.Fatalf("seed table: %v", err)
	}
	guest, err := postgres.NewGuestRepo(testPool).Create(ctx,
		domain.NewGuest("Ivan", "+79001234567"))
	if err != nil {
		t.Fatalf("seed guest: %v", err)
	}
	return tbl.ID, guest.ID
}

func newBooking(tableID, guestID int64, start time.Time) domain.Reservation {
	return domain.Reservation{
		TableID:      tableID,
		GuestID:      guestID,
		GuestsAmount: 2,
		CreatedAt:    time.Now(),
		StartsAt:     start,
		EndsAt:       start.Add(2 * time.Hour),
		Status:       domain.ReservationStatusActive,
	}
}
