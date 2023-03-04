package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"

	"github.com/mustan989/wallet/app/internal/repository/postgres"
	"github.com/mustan989/wallet/model"
	"github.com/mustan989/wallet/repository"
)

func stringp(s string) *string     { return &s }
func boolp(b bool) *bool           { return &b }
func timep(t time.Time) *time.Time { return &t }

func dataToReturnRows(data ...*model.Wallet) *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id", "name", "description", "currency", "amount", "personal", "created_at", "updated_at", "deleted_at",
	}).AddRows(dataToRows(data)...)
}

func dataToRows(data []*model.Wallet) (rows [][]any) {
	for _, wallet := range data {
		rows = append(
			rows, dataToRow(wallet),
		)
	}
	return
}

func dataToRow(data *model.Wallet) []any {
	return []any{
		data.ID, data.Name, data.Description, data.Currency, data.Amount, data.Personal, data.CreatedAt, data.UpdatedAt, data.DeletedAt,
	}
}

var connErr = errors.New("connection error")

func TestWallet_CountAll(t *testing.T) {
	subtests := [...]struct {
		name   string
		expect uint64
		err    error
	}{
		{"0", 0, nil},
		{"1", 1, nil},
		{"connErr", 0, connErr},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	wallet := postgres.NewWallet(pool)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			expectQuery := pool.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM wallets"))
			if subtest.err != nil {
				expectQuery.WillReturnError(subtest.err)
			} else {
				expectQuery.WillReturnRows(
					pgxmock.NewRows([]string{"count"}).AddRow(subtest.expect),
				)
			}

			count, err := wallet.CountAll(context.Background(), &model.WalletFilter{})
			require.Equal(t, err, subtest.err)
			require.Equal(t, subtest.expect, count)
		})
	}
}

func TestWallet_FindAll(t *testing.T) {
	subtests := [...]struct {
		name   string
		expect []*model.Wallet
		err    error
	}{
		{"empty", []*model.Wallet{}, nil},
		{"3 items", []*model.Wallet{
			{ID: 15, Name: "name", Currency: "KZT", Amount: 1011, Personal: true},
			{ID: 17, Name: "name", Currency: "GBP", Amount: 10110001, Personal: false},
			{ID: 128, Name: "name", Currency: "USD", Amount: 667, Personal: true},
		}, nil},
		{"connErr", nil, connErr},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	wallet := postgres.NewWallet(pool)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			expectQuery := pool.ExpectQuery("SELECT (.+) FROM wallets")
			if subtest.err != nil {
				expectQuery.WillReturnError(subtest.err)
			} else {
				expectQuery.WillReturnRows(dataToReturnRows(subtest.expect...))
			}

			data, err := wallet.FindAll(context.Background(), &model.WalletFilter{})
			require.Equal(t, subtest.err, err)
			require.Equal(t, subtest.expect, data)
		})
	}
}

func TestWallet_FindByID(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  uint64
		expect *model.Wallet
		err    error
	}{
		{"success", 15, &model.Wallet{ID: 15, Name: "name", Currency: "KZT", Amount: 1011, Personal: true}, nil},
		{"errWalletNotFound", 15, nil, repository.ErrWalletNotFound},
		{"connErr", 15, nil, connErr},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	wallet := postgres.NewWallet(pool)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			expectQuery := pool.ExpectQuery("SELECT (.+) FROM wallets").WithArgs(subtest.input)
			switch subtest.err {
			case nil:
				expectQuery.WillReturnRows(dataToReturnRows(subtest.expect))
			case repository.ErrWalletNotFound:
				expectQuery.WillReturnError(pgx.ErrNoRows)
			default:
				expectQuery.WillReturnError(subtest.err)
			}

			data, err := wallet.FindByID(context.Background(), subtest.input)
			require.Equal(t, subtest.err, err)
			require.Equal(t, subtest.expect, data)
		})
	}
}

func TestWallet_Create(t *testing.T) {
	subtests := [...]struct {
		name          string
		input, expect *model.Wallet
		err           error
	}{
		{"success", &model.Wallet{
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
			Personal:    true,
		}, &model.Wallet{
			ID:          10,
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
			Personal:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}, nil},
		{"conflict", &model.Wallet{
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
		}, &model.Wallet{
			ID:          10,
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
			Personal:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}, repository.ErrWalletConflict},
		{"connErr", &model.Wallet{
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
			Personal:    true,
		}, nil, connErr},
	}
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	wallet := postgres.NewWallet(pool)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			data := subtest.input
			expect := subtest.expect

			expectQuery := pool.ExpectQuery("INSERT INTO wallets (.+)").WithArgs(
				data.Name, data.Description, data.Currency, data.Amount, data.Personal,
			)

			switch subtest.err {
			case nil:
				expectQuery.WillReturnRows(dataToReturnRows(expect))
			case repository.ErrWalletNotFound:
				expectQuery.WillReturnError(pgx.ErrNoRows)
			default:
				expectQuery.WillReturnError(subtest.err)
			}

			err := wallet.Create(context.Background(), data)
			require.Equal(t, subtest.err, err)

			if subtest.err != nil {
				return
			}

			require.Equal(t, expect, data)
		})
	}
}

func TestWallet_Update(t *testing.T) {
	now := time.Now()

	subtests := [...]struct {
		name          string
		input, expect *model.Wallet
		err           error
	}{
		{"success", &model.Wallet{
			ID:          10,
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
			Personal:    true,
		}, &model.Wallet{
			ID:          10,
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
			Personal:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}, nil},
		{"delete", &model.Wallet{
			ID:          10,
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
			Personal:    true,
			DeletedAt:   timep(now),
		}, &model.Wallet{
			ID:          10,
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
			Personal:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   timep(now),
		}, nil},
		{"conflict", &model.Wallet{
			ID:          10,
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
		}, &model.Wallet{
			ID:          10,
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
			Personal:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}, repository.ErrWalletConflict},
		{"connErr", &model.Wallet{
			ID:          10,
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      1010,
			Personal:    true,
		}, nil, connErr},
	}
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	wallet := postgres.NewWallet(pool)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			data := subtest.input
			expect := subtest.expect

			expectQuery := pool.ExpectQuery("UPDATE wallets").WithArgs(
				data.Name, data.Description, data.Currency, data.Amount, data.Personal, data.DeletedAt, data.ID,
			)

			switch subtest.err {
			case nil:
				expectQuery.WillReturnRows(dataToReturnRows(expect))
			case repository.ErrWalletNotFound:
				expectQuery.WillReturnError(pgx.ErrNoRows)
			default:
				expectQuery.WillReturnError(subtest.err)
			}

			err := wallet.Update(context.Background(), data)
			require.Equal(t, subtest.err, err)

			if subtest.err != nil {
				return
			}

			require.Equal(t, expect, data)
		})
	}
}

func TestWallet_DeleteByID(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  uint64
		expect *model.Wallet
		err    error
	}{
		{"success", 15, &model.Wallet{ID: 15, Name: "name", Currency: "KZT", Amount: 1011, Personal: true}, nil},
		{"errWalletNotFound", 15, nil, repository.ErrWalletNotFound},
		{"connErr", 15, nil, connErr},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	wallet := postgres.NewWallet(pool)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			expectQuery := pool.ExpectQuery("DELETE FROM wallets").WithArgs(subtest.input)
			switch subtest.err {
			case nil:
				expectQuery.WillReturnRows(dataToReturnRows(subtest.expect))
			case repository.ErrWalletNotFound:
				expectQuery.WillReturnError(pgx.ErrNoRows)
			default:
				expectQuery.WillReturnError(subtest.err)
			}

			data, err := wallet.DeleteByID(context.Background(), subtest.input)
			require.Equal(t, subtest.err, err)

			if subtest.err != nil {
				return
			}

			require.Equal(t, subtest.expect, data)
		})
	}
}
