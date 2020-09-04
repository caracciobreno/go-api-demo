package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AccountRepository interface {
	FindUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	ListTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]Transaction, error)
	CreateTransaction(ctx context.Context, transaction *Transaction) error

	FindAndLockUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	UpdateUserBalance(ctx context.Context, userID uuid.UUID, newBalance float64) error

	WithTx(context.Context, func(repository AccountRepository) error) error
}

type Account struct {
	repository AccountRepository
}

func NewAccount(repository AccountRepository) *Account {
	return &Account{repository: repository}
}

func (service *Account) CreateTransaction(ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID, amount float64) (*Transaction, error) {

	var transaction *Transaction
	err := service.repository.WithTx(ctx, func(txRepo AccountRepository) error {
		sourceUser, err := service.repository.FindAndLockUserByID(ctx, sourceUserID)
		if err != nil {
			return err
		}

		targetUser, err := service.repository.FindAndLockUserByID(ctx, targetUserID)
		if err != nil {
			return err
		}

		if sourceUser.Balance < amount {
			return fmt.Errorf("insufficient balance for the transaction")
		}

		sourceUser.Balance -= amount
		if err := txRepo.UpdateUserBalance(ctx, sourceUser.ID, sourceUser.Balance); err != nil {
			return err
		}

		targetUser.Balance += amount
		if err := txRepo.UpdateUserBalance(ctx, targetUser.ID, targetUser.Balance); err != nil {
			return err
		}

		transaction = &Transaction{
			ID:           uuid.New(),
			SourceUserID: sourceUserID,
			TargetUserID: targetUserID,
			Amount:       amount,
			CreatedAt:    time.Now(),
		}

		return txRepo.CreateTransaction(ctx, transaction)
	})

	return transaction, err
}

func (service *Account) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	if userID == uuid.Nil {
		return 0, fmt.Errorf("userID not provided")
	}

	user, err := service.repository.FindUserByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return user.Balance, nil
}

func (service *Account) ListTransactions(ctx context.Context, userID uuid.UUID) ([]Transaction, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("userID not provided")
	}

	return service.repository.ListTransactionsByUserID(ctx, userID)
}
