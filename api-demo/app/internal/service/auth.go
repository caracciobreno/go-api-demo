package service

import (
	"context"
)

// AuthenticationRepository defines a repository that is able to fetch and authenticate users
type AuthenticationRepository interface {

	// FindUserByCredentials looks up in the DB for a user that matches the userName and password combination
	FindUserByCredentials(ctx context.Context, userName string, password string) (*User, error)
}

// Authentication provides implementation of Authentication
type Authentication struct {
	repository AuthenticationRepository
}

func NewAuthentication(repository AuthenticationRepository) *Authentication {
	return &Authentication{repository: repository}
}

func (a *Authentication) Authenticate(ctx context.Context, userName string, password string) (*User, error) {
	return a.repository.FindUserByCredentials(ctx, userName, password)
}
