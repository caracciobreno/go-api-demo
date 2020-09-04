package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AccountRepository defines features that should be provided to the service regarding storage
type AccountRepository interface {

	// FindUserByID looks up for an User with the given ID
	FindUserByID(ctx context.Context, userID uuid.UUID) (*User, error)

	// ListTransactionsByUserID lists all the transactions that a given user was the Source
	ListTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]Transaction, error)

	// CreateTransaction creates a transaction between 2 users
	CreateTransaction(ctx context.Context, transaction *Transaction) error

	// FindAndLockUserByID looks up for an User with the given ID and locks it, not allowing other processes to observe
	// this user while the transaction is not finished
	FindAndLockUserByID(ctx context.Context, userID uuid.UUID) (*User, error)

	// UpdateUserBalance updates the user balance to the given amount
	UpdateUserBalance(ctx context.Context, userID uuid.UUID, newBalance float64) error

	// WithTx starts a transactioned version of the repository that'll be either commited if no errors are returned or
	// rolled back
	WithTx(context.Context, func(repository AccountRepository) error) error
}

// Account provides services related to account
type Account struct {
	repository AccountRepository
}

func NewAccount(repository AccountRepository) *Account {
	return &Account{repository: repository}
}

func (service *Account) CreateTransaction(ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID,
	amount float64) (*Transaction, error) {

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

		if sourceUser.ID == targetUser.ID {
			return errors.New("the target user should be different than the source user")
		}

		if sourceUser.Balance < amount {
			return errors.New("insufficient balance for the transaction")
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

	transactions, err := service.repository.ListTransactionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if transactions == nil {
		transactions = []Transaction{}
	}

	return transactions, nil
}
