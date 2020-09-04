package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"api-demo/app/internal/service"
	customhttp "api-demo/pkg/http"
)

// AuthenticationService defines a service that is able to provide authentication features
type AuthenticationService interface {

	// Authenticate returns a user that matches the userName and password combination
	Authenticate(ctx context.Context, userName string, password string) (*service.User, error)
}

// AuthWrapper wraps a decorated http.HandlerFunc (that receives a user) to a normal one, inspecting the request
// looking for user credentials
type AuthWrapper struct {
	authService AuthenticationService
}

func NewAuthWrapper(authService AuthenticationService) *AuthWrapper {
	return &AuthWrapper{authService: authService}
}

// WithAuth wraps the given function to a normal http.HandlerFunc
func (wrapper *AuthWrapper) WithAuth(f func(w http.ResponseWriter, r *http.Request, user *service.User)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		userName := query.Get("username")
		password := query.Get("password")

		if userName == "" || password == "" {
			customhttp.WriteError(w,
				errors.New("credentials not provided, please provide username and password in the query params"),
				http.StatusUnauthorized)
			return
		}

		user, err := wrapper.authService.Authenticate(r.Context(), userName, password)
		if err != nil {
			customhttp.WriteError(w,
				fmt.Errorf("invalid credentials for user, error: %v", err),
				http.StatusUnauthorized)
			return
		}

		// provides the user to the underlying function
		f(w, r, user)
	}
}
