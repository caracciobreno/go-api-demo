package service

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	UserName string
	Password string
	Balance  float64
}

type Transaction struct {
	ID           uuid.UUID
	SourceUserID uuid.UUID
	TargetUserID uuid.UUID
	Amount       float64
	CreatedAt    time.Time
}
