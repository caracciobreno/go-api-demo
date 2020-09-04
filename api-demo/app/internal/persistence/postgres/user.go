package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"api-demo/app/internal/service"
	"api-demo/pkg/pqutil"
)

const userFields = `id, username, password, balance`

func scanUser(scanner pqutil.Scanner) (*service.User, error) {
	var out service.User
	err := scanner.Scan(&out.ID, &out.UserName, &out.Password, &out.Balance)
	if err == sql.ErrNoRows {
		return nil, errors.New("no such user")
	}
	if err != nil {
		return nil, fmt.Errorf("unexpected error scanning user: %v", err)
	}
	return &out, nil
}
