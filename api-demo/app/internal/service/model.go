package service

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"user_name"`
	Password string    `json:"password"`
	Balance  float64   `json:"balance"`
}

type Transaction struct {
	ID           uuid.UUID `json:"id"`
	SourceUserID uuid.UUID `json:"source_user_id"`
	TargetUserID uuid.UUID `json:"target_user_id"`
	Amount       float64   `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
}
