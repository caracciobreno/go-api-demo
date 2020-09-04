package httpapi

import (
	"context"
	"encoding/json"
	"errors"
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
	service AccountService
}

func NewAccount(service AccountService) *Account {
	return &Account{service: service}
}

func (d *Account) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/{userID}", d.getBalance).Methods(http.MethodGet)
	router.HandleFunc("/{userID}/transactions", d.listTransactions).Methods(http.MethodGet)
	router.HandleFunc("/{userID}/transactions", d.createTransaction).Methods(http.MethodPost)
}

func (d *Account) getBalance(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		customhttp.WriteError(w, errors.New("invalid userID"), http.StatusBadRequest)
		return
	}

	balance, err := d.service.GetBalance(r.Context(), userUUID)
	if err != nil {
		customhttp.WriteError(w, err, http.StatusBadRequest)
		return
	}

	getBalanceResponse := struct {
		UserID  uuid.UUID `json:"user_id"`
		Balance float64   `json:"balance"`
	}{
		userUUID, balance,
	}

	customhttp.WriteJSON(w, getBalanceResponse)
}

func (d *Account) listTransactions(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		customhttp.WriteError(w, errors.New("invalid userID"), http.StatusBadRequest)
		return
	}

	transactions, err := d.service.ListTransactions(r.Context(), userUUID)
	if err != nil {
		customhttp.WriteError(w, err, http.StatusBadRequest)
		return
	}

	listTransactionsResponse := struct {
		UserID       uuid.UUID             `json:"user_id"`
		Transactions []service.Transaction `json:"transactions"`
	}{
		userUUID, transactions,
	}

	customhttp.WriteJSON(w, listTransactionsResponse)
}

func (d *Account) createTransaction(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		customhttp.WriteError(w, errors.New("invalid userID"), http.StatusBadRequest)
		return
	}

	var createTransactionRequest struct {
		TargetUserID uuid.UUID `json:"target_user_id"`
		Amount       float64   `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createTransactionRequest); err != nil {
		customhttp.WriteError(w, err, http.StatusBadRequest)
		return
	}

	transaction, err := d.service.CreateTransaction(r.Context(), userUUID, createTransactionRequest.TargetUserID, createTransactionRequest.Amount)
	if err != nil {
		customhttp.WriteError(w, err, http.StatusBadRequest)
		return
	}

	customhttp.WriteJSON(w, transaction)
}
