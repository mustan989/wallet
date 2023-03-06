package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
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

func getReturnError(err error) error {
	if errors.Is(err, repository.ErrWalletNotFound) {
		return pgx.ErrNoRows
	}
	if errors.Is(err, repository.ErrWalletConflict) {
		return &pgconn.PgError{Code: pgerrcode.IntegrityConstraintViolation}
	}
	return err
}

var connErr = errors.New("connection error")

var rowsAll = []string{
	"id", "name", "description", "currency", "amount", "personal", "created_at", "updated_at", "deleted_at",
}

func TestWallet_CountAll(t *testing.T) {
	subtests := [...]struct {
		name   string
		filter *model.WalletFilter
		args   []any
		expect uint64
	}{
		{"None", &model.WalletFilter{}, nil, 0},
		{"One", &model.WalletFilter{}, nil, 1},
		{"NameLike", &model.WalletFilter{NameLike: "name"}, []any{"%name%"}, 1},
		{"DescriptionLike", &model.WalletFilter{DescriptionLike: "desc"}, []any{"%desc%"}, 1},
		{"Currency", &model.WalletFilter{Currency: "KZT"}, []any{"KZT"}, 1},
		{"Personal", &model.WalletFilter{Personal: boolp(true)}, []any{boolp(true)}, 1},
		{"Currency Personal", &model.WalletFilter{Currency: "KZT", Personal: boolp(true)}, []any{"KZT", boolp(true)}, 1},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			pool.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM wallets")).
				WithArgs(subtest.args...).
				WillReturnRows(pgxmock.NewRows([]string{"count"}).
					AddRow(subtest.expect))

			count, err := wallet.CountAll(context.Background(), subtest.filter)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, count)
		})
	}
}

func TestWallet_CountAllError(t *testing.T) {
	subtests := [...]struct {
		name   string
		filter *model.WalletFilter
		args   []any
		err    error
	}{
		{"None Connection error", &model.WalletFilter{}, nil, connErr},
		{"Currency Connection error", &model.WalletFilter{Currency: "KZT"}, []any{"KZT"}, connErr},
		{"Currency Personal Connection error", &model.WalletFilter{Currency: "KZT", Personal: boolp(true)}, []any{"KZT", boolp(true)}, connErr},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			pool.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM wallets")).
				WithArgs(subtest.args...).
				WillReturnError(getReturnError(subtest.err))

			count, err := wallet.CountAll(context.Background(), subtest.filter)
			require.Zero(t, count)
			require.Equal(t, err, subtest.err)
		})
	}
}

func TestWallet_FindAll(t *testing.T) {
	date := time.Date(1999, 2, 23, 4, 36, 0, 0, time.UTC)

	subtests := [...]struct {
		name   string
		filter *model.WalletFilter
		args   []any
		expect []*model.Wallet
	}{
		{"None", &model.WalletFilter{}, nil, []*model.Wallet{}},
		{"None NameLike", &model.WalletFilter{NameLike: "name"}, []any{"%name%"}, []*model.Wallet{}},
		{"None DescriptionLike", &model.WalletFilter{DescriptionLike: "desc"}, []any{"%desc%"}, []*model.Wallet{}},
		{"None Currency", &model.WalletFilter{Currency: "KZT"}, []any{"KZT"}, []*model.Wallet{}},
		{"None Personal", &model.WalletFilter{Personal: boolp(true)}, []any{boolp(true)}, []*model.Wallet{}},
		{"None Currency Personal", &model.WalletFilter{Currency: "KZT", Personal: boolp(true)}, []any{"KZT", boolp(true)}, []*model.Wallet{}},
		{"Some", &model.WalletFilter{}, nil, []*model.Wallet{
			{1, "name 1", stringp("desc 1"), "KZT", 9999, true, date, date, nil},
			{2, "name 1", stringp("desc 1"), "USD", 9999, false, date, date, nil},
			{3, "name 1", stringp("desc 1"), "EUR", 9999, true, date, date, nil},
			{4, "name 2", stringp("desc 2"), "KZT", 9999, false, date, date, nil},
			{5, "name 2", stringp("desc 2"), "USD", 9999, true, date, date, nil},
			{6, "name 2", stringp("desc 2"), "EUR", 9999, false, date, date, nil},
		}},
		{"Some NameLike", &model.WalletFilter{NameLike: "1"}, []any{"%1%"}, []*model.Wallet{
			{1, "name 1", stringp("desc 1"), "KZT", 9999, true, date, date, nil},
			{2, "name 1", stringp("desc 1"), "USD", 9999, true, date, date, nil},
			{3, "name 1", stringp("desc 1"), "EUR", 9999, false, date, date, nil},
		}},
		{"Some DescriptionLike", &model.WalletFilter{DescriptionLike: "1"}, []any{"%1%"}, []*model.Wallet{
			{1, "name 1", stringp("desc 1"), "KZT", 9999, true, date, date, nil},
			{2, "name 1", stringp("desc 1"), "USD", 9999, true, date, date, nil},
			{3, "name 1", stringp("desc 1"), "EUR", 9999, false, date, date, nil},
		}},
		{"Some Currency", &model.WalletFilter{Currency: "KZT"}, []any{"KZT"}, []*model.Wallet{
			{1, "name 1", stringp("desc 1"), "KZT", 9999, true, date, date, nil},
			{4, "name 2", stringp("desc 2"), "KZT", 9999, true, date, date, nil},
		}},
		{"Some Personal", &model.WalletFilter{Personal: boolp(true)}, []any{boolp(true)}, []*model.Wallet{
			{1, "name 1", stringp("desc 1"), "KZT", 9999, true, date, date, nil},
			{3, "name 1", stringp("desc 1"), "EUR", 9999, true, date, date, nil},
			{5, "name 2", stringp("desc 2"), "USD", 9999, true, date, date, nil},
		}},
		{"Some Currency Personal", &model.WalletFilter{Currency: "KZT", Personal: boolp(true)}, []any{"KZT", boolp(true)}, []*model.Wallet{
			{1, "name 1", stringp("desc 1"), "KZT", 9999, true, date, date, nil},
		}},
		{"Some Limit Offset", &model.WalletFilter{Filter: model.Filter{3, 2}}, nil, []*model.Wallet{
			{3, "name 1", stringp("desc 1"), "EUR", 9999, true, date, date, nil},
			{4, "name 2", stringp("desc 2"), "KZT", 9999, false, date, date, nil},
			{5, "name 2", stringp("desc 2"), "USD", 9999, true, date, date, nil},
		}},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			pool.ExpectQuery("SELECT (.+) FROM wallets").
				WithArgs(subtest.args...).
				WillReturnRows(pgxmock.NewRows(rowsAll).
					AddRows(dataToRows(subtest.expect)...))

			data, err := wallet.FindAll(context.Background(), subtest.filter)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, data)
		})
	}
}

func TestWallet_FindAllError(t *testing.T) {
	subtests := [...]struct {
		name   string
		filter *model.WalletFilter
		args   []any
		err    error
	}{
		{"None Connection error", &model.WalletFilter{}, nil, connErr},
		{"NameLike Connection error", &model.WalletFilter{NameLike: "name"}, []any{"%name%"}, connErr},
		{"DescriptionLike Connection error", &model.WalletFilter{DescriptionLike: "desc"}, []any{"%desc%"}, connErr},
		{"Currency Connection error", &model.WalletFilter{Currency: "KZT"}, []any{"KZT"}, connErr},
		{"Personal Connection error", &model.WalletFilter{Personal: boolp(true)}, []any{boolp(true)}, connErr},
		{"Currency Personal Connection error", &model.WalletFilter{Currency: "KZT", Personal: boolp(true)}, []any{"KZT", boolp(true)}, connErr},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			pool.ExpectQuery("SELECT (.+) FROM wallets").
				WithArgs(subtest.args...).
				WillReturnError(getReturnError(subtest.err))

			data, err := wallet.FindAll(context.Background(), subtest.filter)
			require.Zero(t, data)
			require.Equal(t, subtest.err, err)
		})
	}
}

func TestWallet_FindByID(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  uint64
		expect *model.Wallet
	}{
		{"ID", 1, &model.Wallet{1, "name", stringp("desc"), "KZT", 9999, true, time.Now(), time.Now(), nil}},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			pool.ExpectQuery("SELECT (.+) FROM wallets (.+) LIMIT 1").
				WithArgs(subtest.input).
				WillReturnRows(pgxmock.NewRows(rowsAll).
					AddRow(dataToRow(subtest.expect)...))

			data, err := wallet.FindByID(context.Background(), subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, data)
		})
	}
}

func TestWallet_FindByIDError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input uint64
		err   error
	}{
		{"Not found", 1, repository.ErrWalletNotFound},
		{"Connection error", 1, connErr},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			pool.ExpectQuery("SELECT (.+) FROM wallets (.+) LIMIT 1").
				WithArgs(subtest.input).
				WillReturnError(getReturnError(subtest.err))

			data, err := wallet.FindByID(context.Background(), subtest.input)
			require.Zero(t, data)
			require.Equal(t, subtest.err, err)
		})
	}
}

func TestWallet_Create(t *testing.T) {
	subtests := [...]struct {
		name  string
		input *model.Wallet
	}{
		{"Create", &model.Wallet{
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      9999,
			Personal:    true,
		}},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			data := subtest.input
			now := time.Now()

			pool.ExpectQuery("INSERT INTO wallets (.+)").
				WithArgs(data.Name, data.Description, data.Currency, data.Amount, data.Personal).
				WillReturnRows(pgxmock.NewRows(rowsAll).
					AddRow(uint64(1), data.Name, data.Description, data.Currency, data.Amount, data.Personal, now, now, nil))

			err := wallet.Create(context.Background(), data)
			require.NoError(t, err)

			assert.NotZero(t, data.ID)
			assert.NotZero(t, data.CreatedAt)
			assert.NotZero(t, data.UpdatedAt)
		})
	}
}

func TestWallet_CreateError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input *model.Wallet
		err   error
	}{
		{"Conflict", &model.Wallet{
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      9999,
			Personal:    true,
		}, repository.ErrWalletConflict},
		{"Connection error", &model.Wallet{
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      9999,
			Personal:    true,
		}, connErr},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			data := subtest.input

			pool.ExpectQuery("INSERT INTO wallets (.+)").
				WithArgs(data.Name, data.Description, data.Currency, data.Amount, data.Personal).
				WillReturnError(getReturnError(subtest.err))

			err := wallet.Create(context.Background(), data)
			require.Zero(t, data.ID)
			require.Equal(t, subtest.err, err)
		})
	}
}

func TestWallet_Update(t *testing.T) {
	subtests := [...]struct {
		name  string
		input *model.Wallet
	}{
		{"ID", &model.Wallet{
			ID:          1,
			Name:        "name",
			Description: nil,
			Currency:    "KZT",
			Amount:      9999,
			Personal:    true,
		}},
		{"Delete", &model.Wallet{
			ID:          1,
			Name:        "name",
			Description: nil,
			Currency:    "KZT",
			Amount:      9999,
			Personal:    true,
			DeletedAt:   timep(time.Now()),
		}},
	}
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			data := subtest.input

			now := time.Now()

			pool.ExpectQuery("UPDATE wallets").
				WithArgs(data.Name, data.Description, data.Currency, data.Amount, data.Personal, data.DeletedAt, data.ID).
				WillReturnRows(pgxmock.NewRows(rowsAll).
					AddRow(data.ID, data.Name, data.Description, data.Currency, data.Amount, data.Personal, now.Add(-24*time.Hour), now, data.DeletedAt))

			err := wallet.Update(context.Background(), data)
			require.NoError(t, err)

			assert.Equal(t, now, data.UpdatedAt)
		})
	}
}

func TestWallet_UpdateError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input *model.Wallet
		err   error
	}{
		{"Not found", &model.Wallet{
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      9999,
			Personal:    true,
		}, repository.ErrWalletNotFound},
		{"Conflict", &model.Wallet{
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      9999,
			Personal:    true,
		}, repository.ErrWalletConflict},
		{"Connection error", &model.Wallet{
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      9999,
			Personal:    true,
		}, connErr},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			data := subtest.input

			pool.ExpectQuery("UPDATE wallets").
				WithArgs(data.Name, data.Description, data.Currency, data.Amount, data.Personal, data.DeletedAt, data.ID).
				WillReturnError(getReturnError(subtest.err))

			err := wallet.Update(context.Background(), data)
			require.Zero(t, data.ID)
			require.Equal(t, subtest.err, err)
		})
	}
}

func TestWallet_DeleteByID(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  uint64
		expect *model.Wallet
	}{
		{"ID", 1, &model.Wallet{1, "name", nil, "KZT", 9999, true, time.Now().Add(-24 * time.Hour), time.Now().Add(-10 * time.Minute), timep(time.Now())}},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			pool.ExpectQuery("DELETE FROM wallets").
				WithArgs(subtest.input).
				WillReturnRows(pgxmock.NewRows(rowsAll).
					AddRow(dataToRow(subtest.expect)...))

			data, err := wallet.DeleteByID(context.Background(), subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, data)
		})
	}
}

func TestWallet_DeleteByIDError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input uint64
		err   error
	}{
		{"Not found", 1, repository.ErrWalletNotFound},
		{"Connection error", 1, connErr},
	}

	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			wallet := postgres.NewWallet(pool)

			pool.ExpectQuery("DELETE FROM wallets").
				WithArgs(subtest.input).
				WillReturnError(getReturnError(subtest.err))

			data, err := wallet.DeleteByID(context.Background(), subtest.input)
			require.Zero(t, data)
			require.Equal(t, subtest.err, err)
		})
	}
}
