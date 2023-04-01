package service_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mock_repository "github.com/mustan989/wallet/app/internal/repository/mock"
	. "github.com/mustan989/wallet/app/internal/service"
	"github.com/mustan989/wallet/model"
	"github.com/mustan989/wallet/pkg/logger"
	"github.com/mustan989/wallet/service"
)

var log logger.Logger = logger.NewLogger(logger.WithWriters(map[logger.Severity]io.Writer{
	logger.Debug:   io.Discard,
	logger.Info:    io.Discard,
	logger.Warning: io.Discard,
	logger.Error:   io.Discard,
}))

func stringp(s string) *string     { return &s }
func boolp(b bool) *bool           { return &b }
func timep(t time.Time) *time.Time { return &t }

func TestWallet_Count(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  *service.WalletCountRequest
		expect *service.WalletCountResponse
	}{
		{
			"None",
			&service.WalletCountRequest{Filter: &model.WalletFilter{}},
			&service.WalletCountResponse{Count: 0},
		},
		{
			"One",
			&service.WalletCountRequest{Filter: &model.WalletFilter{}},
			&service.WalletCountResponse{Count: 1},
		},
		{
			"NameLike",
			&service.WalletCountRequest{Filter: &model.WalletFilter{NameLike: "name"}},
			&service.WalletCountResponse{Count: 1},
		},
		{
			"DescriptionLike",
			&service.WalletCountRequest{Filter: &model.WalletFilter{DescriptionLike: "desc"}},
			&service.WalletCountResponse{Count: 1},
		},
		{
			"Currency",
			&service.WalletCountRequest{
				Filter: &model.WalletFilter{Currency: "KZT"},
			},
			&service.WalletCountResponse{Count: 1},
		},
		{
			"Personal",
			&service.WalletCountRequest{Filter: &model.WalletFilter{Personal: boolp(true)}},
			&service.WalletCountResponse{Count: 1},
		},
		{
			"Currency Personal",
			&service.WalletCountRequest{Filter: &model.WalletFilter{Currency: "KZT", Personal: boolp(true)}},
			&service.WalletCountResponse{Count: 1},
		}}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				CountAll(ctx, subtest.input.Filter).
				Return(subtest.expect.Count, nil)

			response, err := svc.Count(ctx, subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, response)
		})
	}
}

func TestWallet_CountError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input *service.WalletCountRequest
		err   error
	}{
		{
			"error",
			&service.WalletCountRequest{Filter: &model.WalletFilter{}},
			errors.New("error"),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				CountAll(ctx, subtest.input.Filter).
				Return(uint64(0), subtest.err)

			response, err := svc.Count(ctx, subtest.input)
			require.Zero(t, response)
			require.Equal(t, subtest.err, err)
		})
	}
}

func TestWallet_GetAll(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  *service.WalletGetAllRequest
		expect *service.WalletGetAllResponse
	}{
		{
			"None",
			&service.WalletGetAllRequest{Filter: &model.WalletFilter{}},
			&service.WalletGetAllResponse{Data: []*model.Wallet{}, Total: 0},
		},
		{
			"One",
			&service.WalletGetAllRequest{Filter: &model.WalletFilter{}},
			&service.WalletGetAllResponse{Data: []*model.Wallet{
				{
					ID:          1,
					Name:        "name 1",
					Description: stringp("desc 1"),
					Currency:    "KZT",
					Amount:      9999,
					Personal:    true,
					CreatedAt:   time.Now().Add(-24 * time.Hour),
					UpdatedAt:   time.Now().Add(-15 * time.Minute),
					DeletedAt:   nil,
				},
			}, Total: 1},
		},
		{
			"NameLike",
			&service.WalletGetAllRequest{Filter: &model.WalletFilter{NameLike: "name"}},
			&service.WalletGetAllResponse{Data: []*model.Wallet{
				{
					ID:          1,
					Name:        "name 1",
					Description: stringp("desc 1"),
					Currency:    "KZT",
					Amount:      9999,
					Personal:    true,
					CreatedAt:   time.Now().Add(-24 * time.Hour),
					UpdatedAt:   time.Now().Add(-15 * time.Minute),
					DeletedAt:   nil,
				},
			}, Total: 1},
		},
		{
			"DescriptionLike",
			&service.WalletGetAllRequest{
				Filter: &model.WalletFilter{DescriptionLike: "desc"},
			},
			&service.WalletGetAllResponse{Data: []*model.Wallet{
				{
					ID:          1,
					Name:        "name 1",
					Description: stringp("desc 1"),
					Currency:    "KZT",
					Amount:      9999,
					Personal:    true,
					CreatedAt:   time.Now().Add(-24 * time.Hour),
					UpdatedAt:   time.Now().Add(-15 * time.Minute),
					DeletedAt:   nil,
				},
			}, Total: 1},
		},
		{
			"Currency",
			&service.WalletGetAllRequest{Filter: &model.WalletFilter{Currency: "KZT"}},
			&service.WalletGetAllResponse{Data: []*model.Wallet{
				{
					ID:          1,
					Name:        "name 1",
					Description: stringp("desc 1"),
					Currency:    "KZT",
					Amount:      9999,
					Personal:    true,
					CreatedAt:   time.Now().Add(-24 * time.Hour),
					UpdatedAt:   time.Now().Add(-15 * time.Minute),
					DeletedAt:   nil,
				},
			}, Total: 1},
		},
		{
			"Personal",
			&service.WalletGetAllRequest{Filter: &model.WalletFilter{Personal: boolp(true)}},
			&service.WalletGetAllResponse{Data: []*model.Wallet{
				{
					ID:          1,
					Name:        "name 1",
					Description: stringp("desc 1"),
					Currency:    "KZT",
					Amount:      9999,
					Personal:    true,
					CreatedAt:   time.Now().Add(-24 * time.Hour),
					UpdatedAt:   time.Now().Add(-15 * time.Minute),
					DeletedAt:   nil,
				},
			}, Total: 1},
		},
		{
			"Currency Personal",
			&service.WalletGetAllRequest{Filter: &model.WalletFilter{Currency: "KZT", Personal: boolp(true)}},
			&service.WalletGetAllResponse{Data: []*model.Wallet{
				{
					ID:          1,
					Name:        "name 1",
					Description: stringp("desc 1"),
					Currency:    "KZT",
					Amount:      9999,
					Personal:    true,
					CreatedAt:   time.Now().Add(-24 * time.Hour),
					UpdatedAt:   time.Now().Add(-15 * time.Minute),
					DeletedAt:   nil,
				},
			}, Total: 1},
		}}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				FindAll(ctx, subtest.input.Filter).
				Return(subtest.expect.Data, nil)

			repo.EXPECT().
				CountAll(ctx, subtest.input.Filter).
				Return(subtest.expect.Total, nil)

			response, err := svc.GetAll(ctx, subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, response)
		})
	}
}

func TestWallet_GetAllError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input *service.WalletGetAllRequest
		err   error
	}{
		{
			"error",
			&service.WalletGetAllRequest{Filter: &model.WalletFilter{}},
			errors.New("error"),
		},
		{
			"count error",
			&service.WalletGetAllRequest{Filter: &model.WalletFilter{}},
			errors.New("error"),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				FindAll(ctx, subtest.input.Filter).
				Return(nil, subtest.err)

			response, err := svc.GetAll(ctx, subtest.input)
			require.Zero(t, response)
			require.Equal(t, subtest.err, err)
		})
	}
}

func TestWallet_GetByID(t *testing.T) {
	subtests := [...]struct {
		name   string
		input  *service.WalletGetByIDRequest
		expect *service.WalletGetByIDResponse
	}{
		{
			"ID", &service.WalletGetByIDRequest{
				ID: 1,
			},
			&service.WalletGetByIDResponse{&model.Wallet{
				ID:          1,
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				DeletedAt:   nil,
			}},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				FindByID(ctx, subtest.input.ID).
				Return(subtest.expect.Data, nil)

			response, err := svc.GetByID(ctx, subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, response)
		})
	}
}

func TestWallet_GetByIDError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input *service.WalletGetByIDRequest
		err   error
	}{
		{
			"error",
			&service.WalletGetByIDRequest{ID: 1},
			errors.New("error"),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				FindByID(ctx, subtest.input.ID).
				Return(nil, subtest.err)

			response, err := svc.GetByID(ctx, subtest.input)
			require.Zero(t, response)
			require.Equal(t, subtest.err, err)
		})
	}
}

func TestWallet_Create(t *testing.T) {
	now := time.Now()

	subtests := [...]struct {
		name   string
		input  *service.WalletCreateRequest
		expect *service.WalletCreateResponse
	}{
		{
			"Create",
			&service.WalletCreateRequest{Data: &model.Wallet{
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
			}},
			&service.WalletCreateResponse{Data: &model.Wallet{
				ID:          1,
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   nil,
			}},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				Create(ctx, subtest.input.Data).
				DoAndReturn(func(_ context.Context, data *model.Wallet) error {
					data.ID = 1
					data.CreatedAt = subtest.expect.Data.CreatedAt
					data.UpdatedAt = subtest.expect.Data.UpdatedAt
					return nil
				})

			response, err := svc.Create(ctx, subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, response)
		})
	}
}

func TestWallet_CreateError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input *service.WalletCreateRequest
		err   error
	}{
		{
			"error",
			&service.WalletCreateRequest{Data: &model.Wallet{
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
			}},
			errors.New("error"),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				Create(ctx, subtest.input.Data).
				Return(subtest.err)

			response, err := svc.Create(ctx, subtest.input)
			require.Zero(t, response)
			require.Equal(t, subtest.err, err)
		})
	}
}

func TestWallet_Update(t *testing.T) {
	now := time.Now()

	subtests := [...]struct {
		name   string
		input  *service.WalletUpdateRequest
		expect *service.WalletUpdateResponse
	}{
		{
			"Update",
			&service.WalletUpdateRequest{Data: &model.Wallet{
				ID:          1,
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
			}},
			&service.WalletUpdateResponse{Data: &model.Wallet{
				ID:          1,
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
				CreatedAt:   now.Add(-24 * time.Hour),
				UpdatedAt:   now,
				DeletedAt:   nil,
			}},
		},
		{
			"Delete",
			&service.WalletUpdateRequest{Data: &model.Wallet{
				ID:          1,
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
				DeletedAt:   timep(now),
			}},
			&service.WalletUpdateResponse{Data: &model.Wallet{
				ID:          1,
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
				CreatedAt:   now.Add(-24 * time.Hour),
				UpdatedAt:   now,
				DeletedAt:   timep(now),
			}},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				Update(ctx, subtest.input.Data).
				DoAndReturn(func(_ context.Context, data *model.Wallet) error {
					data.CreatedAt = subtest.expect.Data.CreatedAt
					data.UpdatedAt = subtest.expect.Data.UpdatedAt
					data.DeletedAt = subtest.expect.Data.DeletedAt
					return nil
				})

			response, err := svc.Update(ctx, subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, response)
		})
	}
}

func TestWallet_UpdateError(t *testing.T) {
	subtests := [...]struct {
		name  string
		input *service.WalletUpdateRequest
		err   error
	}{
		{
			"Update",
			&service.WalletUpdateRequest{Data: &model.Wallet{
				ID:          1,
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
			}},
			errors.New("error"),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				Update(ctx, subtest.input.Data).
				Return(subtest.err)

			response, err := svc.Update(ctx, subtest.input)
			require.Zero(t, response)
			require.Equal(t, subtest.err, err)
		})
	}
}

func TestWallet_DeleteByID(t *testing.T) {
	now := time.Now()

	subtests := [...]struct {
		name   string
		input  *service.WalletDeleteByIDRequest
		expect *service.WalletDeleteByIDResponse
	}{
		{
			"Delete",
			&service.WalletDeleteByIDRequest{ID: 1},
			&service.WalletDeleteByIDResponse{Data: &model.Wallet{
				ID:          1,
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
				CreatedAt:   now.Add(-24 * time.Hour),
				UpdatedAt:   now.Add(-12 * time.Hour),
				DeletedAt:   timep(now.Add(-12 * time.Hour)),
			}},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			repo.EXPECT().
				FindByID(ctx, subtest.input.ID).
				Return(subtest.expect.Data, nil)

			repo.EXPECT().
				DeleteByID(ctx, subtest.input.ID).
				Return(subtest.expect.Data, nil)

			response, err := svc.DeleteByID(ctx, subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, response)
		})
	}
}

func TestWallet_DeleteByIDSoft(t *testing.T) {
	now := time.Now()

	subtests := [...]struct {
		name   string
		input  *service.WalletDeleteByIDRequest
		expect *service.WalletDeleteByIDResponse
	}{
		{
			"Delete",
			&service.WalletDeleteByIDRequest{ID: 1},
			&service.WalletDeleteByIDResponse{Data: &model.Wallet{
				ID:          1,
				Name:        "name",
				Description: stringp("desc"),
				Currency:    "KZT",
				Amount:      9999,
				Personal:    true,
				CreatedAt:   now.Add(-24 * time.Hour),
				UpdatedAt:   now,
				DeletedAt:   timep(now),
			}},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ctx := context.Background()
			ctl := gomock.NewController(t)
			repo := mock_repository.NewMockWallet(ctl)
			svc := NewWallet(repo, WithLogger(log))

			data := subtest.expect.Data
			repo.EXPECT().
				FindByID(ctx, subtest.input.ID).
				DoAndReturn(func(_ context.Context, id uint64) (*model.Wallet, error) {
					data.DeletedAt = nil
					return data, nil
				})

			repo.EXPECT().
				Update(ctx, data).
				DoAndReturn(func(_ context.Context, data *model.Wallet) error {
					data.DeletedAt = subtest.expect.Data.DeletedAt
					return nil
				})

			response, err := svc.DeleteByID(ctx, subtest.input)
			require.NoError(t, err)
			require.Equal(t, subtest.expect, response)
		})
	}
}
