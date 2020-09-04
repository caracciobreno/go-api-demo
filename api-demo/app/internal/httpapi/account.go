package httpapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"api-demo/app/internal/service"
	customhttp "api-demo/pkg/http"
)

// AccountService abstracts the services that should be provided to the HTTP API
type AccountService interface {

	// CreateTransaction creates a transaction to transfer amount from sourceUserID to targetUserID
	CreateTransaction(ctx context.Context, sourceUserID uuid.UUID, targetUserID uuid.UUID, amount float64) (*service.Transaction, error)

	// GetBalance retrieves the balance of the user
	GetBalance(ctx context.Context, userID uuid.UUID) (float64, error)

	// ListTransactions list all the transaction from a certain User
	ListTransactions(ctx context.Context, userID uuid.UUID) ([]service.Transaction, error)
}

type Account struct {
	accountService AccountService
	authWrapper    *AuthWrapper
}

func NewAccount(accountService AccountService, authWrapper *AuthWrapper) *Account {
	return &Account{accountService: accountService, authWrapper: authWrapper}
}

func (d *Account) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/me", d.authWrapper.WithAuth(d.getBalance)).Methods(http.MethodGet)
	router.HandleFunc("/me/transactions", d.authWrapper.WithAuth(d.listTransactions)).Methods(http.MethodGet)
	router.HandleFunc("/me/transactions", d.authWrapper.WithAuth(d.createTransaction)).Methods(http.MethodPost)
}

func (d *Account) getBalance(w http.ResponseWriter, r *http.Request, user *service.User) {

	balance, err := d.accountService.GetBalance(r.Context(), user.ID)
	if err != nil {
		customhttp.WriteError(w, err, http.StatusBadRequest)
		return
	}

	getBalanceResponse := struct {
		UserID  uuid.UUID `json:"user_id"`
		Balance float64   `json:"balance"`
	}{
		user.ID, balance,
	}

	customhttp.WriteJSON(w, getBalanceResponse)
}

func (d *Account) listTransactions(w http.ResponseWriter, r *http.Request, user *service.User) {

	transactions, err := d.accountService.ListTransactions(r.Context(), user.ID)
	if err != nil {
		customhttp.WriteError(w, err, http.StatusBadRequest)
		return
	}

	listTransactionsResponse := struct {
		UserID       uuid.UUID             `json:"user_id"`
		Transactions []service.Transaction `json:"transactions"`
	}{
		user.ID, transactions,
	}

	customhttp.WriteJSON(w, listTransactionsResponse)
}

func (d *Account) createTransaction(w http.ResponseWriter, r *http.Request, user *service.User) {

	var createTransactionRequest struct {
		TargetUserID uuid.UUID `json:"target_user_id"`
		Amount       float64   `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createTransactionRequest); err != nil {
		customhttp.WriteError(w, err, http.StatusBadRequest)
		return
	}

	transaction, err := d.accountService.CreateTransaction(r.Context(), user.ID, createTransactionRequest.TargetUserID,
		createTransactionRequest.Amount)
	if err != nil {
		customhttp.WriteError(w, err, http.StatusBadRequest)
		return
	}

	customhttp.WriteJSON(w, transaction)
}
