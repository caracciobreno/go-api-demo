package httpapi

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			customhttp.WriteError(w,
				errors.New("no authorization provided"),
				http.StatusUnauthorized)
			return
		}

		splitAuthHeader := strings.Split(authHeader, " ")
		if len(splitAuthHeader) != 2 || splitAuthHeader[0] != "Basic" {
			customhttp.WriteError(w,
				errors.New("invalid authorization header provided"),
				http.StatusUnauthorized)
			return
		}

		digest, err := base64.StdEncoding.DecodeString(splitAuthHeader[1])
		if err != nil {
			customhttp.WriteError(w,
				errors.New("failed to decode base64 basic auth content"),
				http.StatusUnauthorized)
			return
		}

		authContent := strings.Split(string(digest), ":")
		if len(authContent) != 2 {
			customhttp.WriteError(w,
				errors.New("invalid format for username:password"),
				http.StatusUnauthorized)
			return
		}

		userName := authContent[0]
		password := authContent[1]

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
