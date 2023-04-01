package service

import (
	"context"

	"github.com/mustan989/wallet/model"
)

type Wallet interface {
	Count(ctx context.Context, request *WalletCountRequest) (*WalletCountResponse, error)
	GetAll(ctx context.Context, request *WalletGetAllRequest) (*WalletGetAllResponse, error)
	GetByID(ctx context.Context, request *WalletGetByIDRequest) (*WalletGetByIDResponse, error)
	Create(ctx context.Context, request *WalletCreateRequest) (*WalletCreateResponse, error)
	Update(ctx context.Context, request *WalletUpdateRequest) (*WalletUpdateResponse, error)
	DeleteByID(ctx context.Context, request *WalletDeleteByIDRequest) (*WalletDeleteByIDResponse, error)
}

type WalletCountRequest struct {
	Filter *model.WalletFilter
}

type WalletCountResponse struct {
	Count uint64
}

type WalletGetAllRequest struct {
	Filter *model.WalletFilter
}

type WalletGetAllResponse struct {
	Data  []*model.Wallet
	Total uint64
}

type WalletGetByIDRequest struct {
	ID uint64
}

type WalletGetByIDResponse struct {
	Data *model.Wallet
}

type WalletCreateRequest struct {
	Data *model.Wallet
}

type WalletCreateResponse struct {
	Data *model.Wallet
}

type WalletUpdateRequest struct {
	Data *model.Wallet
}

type WalletUpdateResponse struct {
	Data *model.Wallet
}

type WalletDeleteByIDRequest struct {
	ID uint64
}

type WalletDeleteByIDResponse struct {
	Data *model.Wallet
}
