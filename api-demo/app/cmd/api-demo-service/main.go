package main

import (
	"context"

	"api-demo/app/internal/httpapi"
	"api-demo/app/internal/persistence/postgres"
	"api-demo/app/internal/service"
	"api-demo/pkg/app"
)

func main() {
	app.New(func(ctx context.Context, resources app.SetupResourcesProvider) error {

		db, err := resources.WithPostgresConnection("postgres")
		if err != nil {
			return err
		}

		accountRepo := postgres.NewAccountRepository(db)
		accountService := service.NewAccount(accountRepo)

		accountAPI := httpapi.NewAccount(accountService)

		resources.WithHTTPAPI(accountAPI)

		return nil
	}).Run()
}
