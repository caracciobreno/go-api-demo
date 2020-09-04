package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"api-demo/app/internal/service"
	"api-demo/pkg/pqutil"
)

const transactionFields = `id, source_user_id, target_user_id, amount, created_at`

func scanTransaction(scanner pqutil.Scanner) (*service.Transaction, error) {
	var out service.Transaction
	err := scanner.Scan(&out.ID, &out.SourceUserID, &out.TargetUserID, &out.Amount, &out.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("no such transaction")
	}
	if err != nil {
		return nil, fmt.Errorf("unexpected error scanning transactions: %v", err)
	}
	return &out, nil
}

func collectTransactions(scanner pqutil.ScannerIter) ([]service.Transaction, error) {
	var transactions []service.Transaction
	for scanner.Next() {
		user, err := scanTransaction(scanner)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *user)
	}
	return transactions, scanner.Err()
}
