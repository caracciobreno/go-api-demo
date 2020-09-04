package httpapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"api-demo/app/internal/service"
)

type AccountService interface {
	// CreateTransaction creates a transaction to transfer amount from sourceUserID to targetUserID
	CreateTransaction(ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID, amount float64) (*service.Transaction, error)

	// GetBalance retrieves the balance of the user
	GetBalance(ctx context.Context, userID uuid.UUID) (float64, error)

	// ListTransactions list all the transaction from a certain User
	ListTransactions(ctx context.Context, userID uuid.UUID) ([]service.Transaction, error)
}

type Account struct {
	service AccountService
}

func NewAccount(service AccountService) *Account {
	return &Account{service: service}
}

func (d *Account) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, err := d.service.CreateTransaction(ctx, uuid.MustParse("126bea59-c9a7-44d0-bcd8-d710aad69676"),
			uuid.MustParse("256bea59-c9a7-44d0-bcd8-d710aad69676"), 400)

		// trans, err := d.service.GetBalance(ctx, uuid.MustParse("256bea59-c9a7-44d0-bcd8-d710aad69676"))
		encoder := json.NewEncoder(w)
		if err != nil {
			err2 := encoder.Encode(err)
			if err2 != nil {
				panic(err2)
			}
			w.WriteHeader(http.StatusBadRequest)

		}

		w.WriteHeader(http.StatusOK)
	})
}
