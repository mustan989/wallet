package service

import (
	"context"
	"time"

	"github.com/mustan989/wallet/pkg/helper"
	"github.com/mustan989/wallet/pkg/logger"
	"github.com/mustan989/wallet/repository"
	"github.com/mustan989/wallet/service"
)

type WalletOption func(w *wallet)

func WithLogger(log logger.Logger) WalletOption { return func(w *wallet) { w.log = log } }

func NewWallet(repo repository.Wallet, options ...WalletOption) service.Wallet {
	w := &wallet{
		log:  logger.Default(),
		repo: repo,
	}

	for _, option := range options {
		option(w)
	}

	return w
}

type wallet struct {
	log logger.Logger

	repo repository.Wallet
}

func (w *wallet) Count(ctx context.Context, request *service.WalletCountRequest) (*service.WalletCountResponse, error) {
	count, err := w.repo.CountAll(ctx, request.Filter)
	if err != nil {
		w.log.Errorf("Error getting wallet count: %s", err)
		return nil, err
	}
	return &service.WalletCountResponse{Count: count}, nil
}

func (w *wallet) GetAll(ctx context.Context, request *service.WalletGetAllRequest) (*service.WalletGetAllResponse, error) {
	data, err := w.repo.FindAll(ctx, request.Filter)
	if err != nil {
		w.log.Errorf("Error getting wallets: %s", err)
		return nil, err
	}

	count, err := w.Count(ctx, &service.WalletCountRequest{Filter: request.Filter})
	if err != nil {
		return nil, err
	}

	return &service.WalletGetAllResponse{
		Data:  data,
		Total: count.Count,
	}, nil
}

func (w *wallet) GetByID(ctx context.Context, request *service.WalletGetByIDRequest) (*service.WalletGetByIDResponse, error) {
	data, err := w.repo.FindByID(ctx, request.ID)
	if err != nil {
		w.log.Errorf("Error getting wallet by id %d: %s", request.ID, err)
		return nil, err
	}
	return &service.WalletGetByIDResponse{Data: data}, nil
}

func (w *wallet) Create(ctx context.Context, request *service.WalletCreateRequest) (*service.WalletCreateResponse, error) {
	if err := w.repo.Create(ctx, request.Data); err != nil {
		w.log.Errorf("Error creating wallet: %s", err)
		return nil, err
	}
	return &service.WalletCreateResponse{Data: request.Data}, nil
}

func (w *wallet) Update(ctx context.Context, request *service.WalletUpdateRequest) (*service.WalletUpdateResponse, error) {
	if err := w.repo.Update(ctx, request.Data); err != nil {
		w.log.Errorf("Error updating wallet: %s", err)
		return nil, err
	}
	return &service.WalletUpdateResponse{Data: request.Data}, nil
}

func (w *wallet) DeleteByID(ctx context.Context, request *service.WalletDeleteByIDRequest) (*service.WalletDeleteByIDResponse, error) {
	data, err := w.repo.FindByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	if data.DeletedAt == nil {
		data.DeletedAt = helper.Timep(time.Now())
		if err = w.repo.Update(ctx, data); err != nil {
			w.log.Errorf("Error marking wallet as deleted: %s", err)
			return nil, err
		}
		return &service.WalletDeleteByIDResponse{Data: data}, nil
	}

	deleted, err := w.repo.DeleteByID(ctx, request.ID)
	if err != nil {
		w.log.Errorf("Error deleting wallet: %s", err)
		return nil, err
	}
	return &service.WalletDeleteByIDResponse{Data: deleted}, nil
}
