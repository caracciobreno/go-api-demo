package service_test

import (
	"context"

	"github.com/google/uuid"

	"api-demo/app/internal/service"
)

type accountRepositoryMock struct {
	FindUserByIDFunc             func(ctx context.Context, userID uuid.UUID) (*service.User, error)
	ListTransactionsByUserIDFunc func(ctx context.Context, userID uuid.UUID) ([]service.Transaction, error)
	CreateTransactionFunc        func(ctx context.Context, transaction *service.Transaction) error
	FindAndLockUserByIDFunc      func(ctx context.Context, userID uuid.UUID) (*service.User, error)
	UpdateUserBalanceFunc        func(ctx context.Context, userID uuid.UUID, newBalance float64) error
}

func newAccountRepositoryMock() *accountRepositoryMock {
	mock := &accountRepositoryMock{
		FindUserByIDFunc: func(context.Context, uuid.UUID) (*service.User, error) {
			return nil, nil
		},
		ListTransactionsByUserIDFunc: func(context.Context, uuid.UUID) ([]service.Transaction, error) {
			return nil, nil
		},
		CreateTransactionFunc: func(context.Context, *service.Transaction) error {
			return nil
		},
		FindAndLockUserByIDFunc: func(context.Context, uuid.UUID) (*service.User, error) {
			return nil, nil
		},
		UpdateUserBalanceFunc: func(context.Context, uuid.UUID, float64) error {
			return nil
		},
	}

	return mock
}

func (a *accountRepositoryMock) FindUserByID(ctx context.Context, userID uuid.UUID) (*service.User, error) {
	return a.FindUserByIDFunc(ctx, userID)
}

func (a *accountRepositoryMock) ListTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]service.Transaction, error) {
	return a.ListTransactionsByUserIDFunc(ctx, userID)
}

func (a *accountRepositoryMock) CreateTransaction(ctx context.Context, transaction *service.Transaction) error {
	return a.CreateTransactionFunc(ctx, transaction)
}

func (a *accountRepositoryMock) FindAndLockUserByID(ctx context.Context, userID uuid.UUID) (*service.User, error) {
	return a.FindAndLockUserByIDFunc(ctx, userID)
}

func (a *accountRepositoryMock) UpdateUserBalance(ctx context.Context, userID uuid.UUID, newBalance float64) error {
	return a.UpdateUserBalanceFunc(ctx, userID, newBalance)
}

func (a *accountRepositoryMock) WithTx(ctx context.Context, f func(repository service.AccountRepository) error) error {
	return f(a)
}
