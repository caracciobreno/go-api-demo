package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"api-demo/app/internal/service"
	"api-demo/pkg/pqutil"
)

// AccountRepository is a persistence repository that uses Postgres as its DB
type AccountRepository struct {
	queryer pqutil.Queryer
	txer    pqutil.Transactioner
}

// NewAccountRepository creates a postgres repository for accounts
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{queryer: db, txer: db}
}

func (repo *AccountRepository) FindUserByID(ctx context.Context, userID uuid.UUID) (*service.User, error) {
	const query = `SELECT ` + userFields + ` FROM users WHERE id = $1`
	return scanUser(repo.queryer.QueryRowContext(ctx, query, userID))
}

func (repo *AccountRepository) ListTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]service.Transaction, error) {
	const query = `SELECT ` + transactionFields + ` FROM transactions WHERE source_user_id = $1`

	rows, err := repo.queryer.QueryContext(ctx, query,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("unexpected error listing transactions: %v", err)
	}

	defer rows.Close()
	transactions, err := collectTransactions(rows)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *AccountRepository) CreateTransaction(ctx context.Context, transaction *service.Transaction) error {

	const insertQuery = `INSERT INTO transactions (` + transactionFields + `) VALUES ($1, $2, $3, $4, $5)`

	_, err := repo.queryer.ExecContext(ctx, insertQuery,
		transaction.ID,
		transaction.SourceUserID,
		transaction.TargetUserID,
		transaction.Amount,
		transaction.CreatedAt,
	)

	return err
}

func (repo *AccountRepository) FindAndLockUserByID(ctx context.Context, userID uuid.UUID) (*service.User, error) {
	const query = `SELECT ` + userFields + ` FROM users WHERE id = $1 FOR UPDATE`
	return scanUser(repo.queryer.QueryRowContext(ctx, query, userID))
}

func (repo *AccountRepository) UpdateUserBalance(ctx context.Context, userID uuid.UUID, newBalance float64) error {

	const updateQuery = `UPDATE users SET balance = $2 WHERE ID = $1`

	_, err := repo.queryer.ExecContext(ctx, updateQuery,
		userID,
		newBalance,
	)

	return err
}

func (repo *AccountRepository) WithTx(ctx context.Context, transactionedFunction func(repository service.AccountRepository) error) error {
	tx, err := repo.txer.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	// creates a new version of the repo but transactioned
	err = transactionedFunction(&AccountRepository{queryer: tx, txer: nil})
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			panic(fmt.Sprintf("failed to rollback transaction: %v", txErr))
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		panic(fmt.Sprintf("failed to commit transaction: %v", err))
	}

	return nil
}

func (repo *AccountRepository) FindUserByCredentials(ctx context.Context, userName string, password string) (*service.User, error) {
	const query = `SELECT ` + userFields + ` FROM users WHERE username = $1 AND password = $2`
	return scanUser(repo.queryer.QueryRowContext(ctx, query, userName, password))
}
