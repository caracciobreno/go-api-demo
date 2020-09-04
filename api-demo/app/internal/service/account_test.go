package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"api-demo/app/internal/service"
)

func TestAccount_GetBalance(t *testing.T) {

	ctx := context.Background()

	successCheck := func(t *testing.T, user *service.User, balance float64, err error) {
		require.NoError(t, err)
		require.Equal(t, user.Balance, balance)
	}

	failCheck := func(t *testing.T, user *service.User, balance float64, err error) {
		require.Error(t, err)
		require.Zero(t, balance)
	}

	tests := map[string]struct {
		mutateUser    func(*service.User)
		mutateMock    func(*accountRepositoryMock)
		checkFunction func(*testing.T, *service.User, float64, error)
	}{
		"should succeed to get an user's balance": {
			checkFunction: successCheck,
		},
		"should return an error when the DB layer returns an error": {
			mutateMock: func(mock *accountRepositoryMock) {
				mock.FindUserByIDFunc = func(ctx context.Context, userID uuid.UUID) (*service.User, error) {
					return nil, errors.New("not found")
				}
			},
			checkFunction: failCheck,
		},
		"should return an error when the userID is nil": {
			mutateUser: func(user *service.User) {
				user.ID = uuid.Nil
			},
			checkFunction: failCheck,
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {

			repo := newAccountRepositoryMock()

			userID := uuid.New()
			user := &service.User{
				ID:       userID,
				UserName: userID.String(),
				Password: userID.String(),
				Balance:  100,
			}

			if test.mutateUser != nil {
				test.mutateUser(user)
			}

			repo.FindUserByIDFunc = func(ctx context.Context, userID uuid.UUID) (*service.User, error) {
				return user, nil
			}

			if test.mutateMock != nil {
				test.mutateMock(repo)
			}

			accountService := service.NewAccount(repo)

			balance, err := accountService.GetBalance(ctx, user.ID)
			test.checkFunction(t, user, balance, err)

		})
	}
}

func TestAccount_ListTransactions(t *testing.T) {

	ctx := context.Background()

	successCheck := func(t *testing.T, user *service.User, transactions []service.Transaction, err error) {
		require.NoError(t, err)
		require.NotEmpty(t, transactions)

		for _, transaction := range transactions {
			require.Equal(t, user.ID, transaction.SourceUserID)
			require.NotEqual(t, user.ID, transaction.TargetUserID)
		}
	}

	failCheck := func(t *testing.T, user *service.User, transactions []service.Transaction, err error) {
		require.Error(t, err)
		require.Nil(t, transactions)
	}

	tests := map[string]struct {
		mutateUser    func(*service.User)
		mutateMock    func(*accountRepositoryMock)
		checkFunction func(*testing.T, *service.User, []service.Transaction, error)
	}{
		"should succeed to get an user's transaction": {
			checkFunction: successCheck,
		},
		"should return an error when the DB layer returns an error": {
			mutateMock: func(mock *accountRepositoryMock) {
				mock.ListTransactionsByUserIDFunc = func(ctx context.Context, userID uuid.UUID) ([]service.Transaction, error) {
					return nil, errors.New("not found")
				}
			},
			checkFunction: failCheck,
		},
		"should return an error when the userID is nil": {
			mutateUser: func(user *service.User) {
				user.ID = uuid.Nil
			},
			checkFunction: failCheck,
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {

			repo := newAccountRepositoryMock()

			userID := uuid.New()
			user := &service.User{
				ID:       userID,
				UserName: userID.String(),
				Password: userID.String(),
				Balance:  100,
			}

			if test.mutateUser != nil {
				test.mutateUser(user)
			}

			repo.ListTransactionsByUserIDFunc = func(ctx context.Context, userID uuid.UUID) ([]service.Transaction, error) {
				return []service.Transaction{
					{
						ID:           uuid.New(),
						SourceUserID: userID,
						TargetUserID: uuid.New(),
						Amount:       10,
						CreatedAt:    time.Now(),
					},
				}, nil
			}

			if test.mutateMock != nil {
				test.mutateMock(repo)
			}

			accountService := service.NewAccount(repo)

			transactions, err := accountService.ListTransactions(ctx, user.ID)
			test.checkFunction(t, user, transactions, err)

		})
	}
}

func TestAccount_CreateTransaction(t *testing.T) {

	ctx := context.Background()

	successCheck := func(t *testing.T, sourceUser *service.User, targetUser *service.User, sourceOriginalBalance float64,
		targetOriginalBalance float64, transaction *service.Transaction, err error) {
		require.NoError(t, err)
		require.NotNil(t, transaction)

		// guarantees that the balance is correct calculated and that the amount is also correct on the transaction
		require.Equal(t, transaction.Amount, sourceOriginalBalance-sourceUser.Balance, "invalid new balance for source user")
		require.Equal(t, transaction.Amount, targetUser.Balance-targetOriginalBalance, "invalid new balance for target user")

		require.Equal(t, sourceUser.ID, transaction.SourceUserID)
		require.Equal(t, targetUser.ID, transaction.TargetUserID)
	}

	failCheck := func(t *testing.T, sourceUser *service.User, targetUser *service.User, sourceOriginalBalance float64,
		targetOriginalBalance float64, transaction *service.Transaction, err error) {
		require.Error(t, err)
		require.Nil(t, transaction)
	}

	tests := map[string]struct {
		mutateSourceUser func(*service.User)
		mutateTargetUser func(*service.User)
		mutateMock       func(*accountRepositoryMock)
		amount           float64
		checkFunction    func(*testing.T, *service.User, *service.User, float64, float64, *service.Transaction, error)
	}{
		"should succeed creating a transaction when an user has more than the amount as balance": {
			amount: 5,
			mutateSourceUser: func(user *service.User) {
				user.Balance = 10
			},
			checkFunction: successCheck,
		},
		"should succeed creating a transaction when an user has more the same  amount as balance": {
			amount: 10,
			mutateSourceUser: func(user *service.User) {
				user.Balance = 10
			},
			checkFunction: successCheck,
		},
		"should return an error when the DB doesn't find an user": {
			mutateMock: func(mock *accountRepositoryMock) {
				mock.FindAndLockUserByIDFunc = func(ctx context.Context, userID uuid.UUID) (*service.User, error) {
					return nil, errors.New("couldn't find user")
				}
			},
			checkFunction: failCheck,
		},
		"should return an error when the source user have insufficient funds": {
			amount: 100,
			mutateSourceUser: func(user *service.User) {
				user.Balance = 5
			},
			checkFunction: func(t *testing.T, user *service.User, user2 *service.User, f float64, f2 float64, transaction *service.Transaction, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "insufficient balance")
			},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {

			repo := newAccountRepositoryMock()

			sourceUserID := uuid.New()
			sourceUser := &service.User{
				ID:       sourceUserID,
				UserName: sourceUserID.String(),
				Password: sourceUserID.String(),
				Balance:  100,
			}

			targetUserID := uuid.New()
			targetUser := &service.User{
				ID:       targetUserID,
				UserName: targetUserID.String(),
				Password: targetUserID.String(),
				Balance:  100,
			}

			if test.mutateSourceUser != nil {
				test.mutateSourceUser(sourceUser)
			}

			if test.mutateTargetUser != nil {
				test.mutateTargetUser(targetUser)
			}

			sourceUserBalance := sourceUser.Balance
			targetUserBalance := targetUser.Balance

			users := map[uuid.UUID]*service.User{
				sourceUserID: sourceUser,
				targetUserID: targetUser,
			}

			repo.FindAndLockUserByIDFunc = func(ctx context.Context, userID uuid.UUID) (*service.User, error) {
				return users[userID], nil
			}

			if test.mutateMock != nil {
				test.mutateMock(repo)
			}

			accountService := service.NewAccount(repo)

			transaction, err := accountService.CreateTransaction(ctx, sourceUserID, targetUserID, test.amount)
			test.checkFunction(t, sourceUser, targetUser, sourceUserBalance, targetUserBalance, transaction, err)
		})
	}
}
