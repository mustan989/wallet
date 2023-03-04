package repository

import (
	"context"
	"errors"

	"github.com/mustan989/wallet/model"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
	ErrWalletConflict = errors.New("wallet already exists")
)

// Wallet repository interface
type Wallet interface {
	CountAll(ctx context.Context, filter *model.WalletFilter) (count uint64, err error)
	FindAll(ctx context.Context, filter *model.WalletFilter) (data []*model.Wallet, err error)
	FindByID(ctx context.Context, id uint64) (data *model.Wallet, err error)
	Create(ctx context.Context, data *model.Wallet) error
	Update(ctx context.Context, data *model.Wallet) error
	DeleteByID(ctx context.Context, id uint64) (deleted *model.Wallet, err error)
}
