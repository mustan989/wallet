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
	logger2 "github.com/mustan989/wallet/pkg/logger"
	"github.com/mustan989/wallet/service"
)

var log logger2.Logger = logger2.NewLogger(logger2.WithWriters(map[logger2.Severity]io.Writer{
	logger2.Debug:   io.Discard,
	logger2.Info:    io.Discard,
	logger2.Warning: io.Discard,
	logger2.Error:   io.Discard,
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
		{"ID", &service.WalletGetByIDRequest{
			ID: 1,
		}, &service.WalletGetByIDResponse{&model.Wallet{
			ID:          1,
			Name:        "name",
			Description: stringp("desc"),
			Currency:    "KZT",
			Amount:      9999,
			Personal:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}}},
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
}
