package app

import (
	"context"
	"database/sql"
	"fmt"
	gohttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	// register the pg driver
	_ "github.com/lib/pq"

	"api-demo/pkg/http"
	"api-demo/pkg/log"
)

// App is an abstraction of a service, or standalone application running whose lifecycle is automatically managed.
// The idea of this abstraction is to allow the creation of integration tests reusing the same Setup that is provided
// on a real environment, just changing the implementation of an App and a SetupResourcesProvider.
type App interface {

	// Run starts the app, starting all its subcomponents
	Run()
}

// SetupResourcesProvider provides, during the Setup of the App, resources for the underlying app
type SetupResourcesProvider interface {

	// WithHTTPAPI uses registers the given API on a HTTP server
	WithHTTPAPI(api http.API)

	WithPostgresConnection(db string) (*sql.DB, error)
}

// Shutdowner defines something that can shutdown
type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

type SetupFunc func(context.Context, SetupResourcesProvider) error

// StandardApp is a real implementation of an App
type StandardApp struct {
	setupFunc  SetupFunc
	router     *mux.Router
	apis       []http.API
	toShutdown []Shutdowner
}

// New creates a Standard App ready to be configured using the setupFunc
func New(setupFunc SetupFunc) *StandardApp {
	return &StandardApp{setupFunc: setupFunc, router: basicRouter()}
}

func (app *StandardApp) Run() {

	logger := log.New(
		log.WithFormatter(log.DefaultFormatter),
		log.WithLevel(logrus.InfoLevel),
		log.WithOutput(os.Stderr))

	// sets the logger in the context, the intention is to reuse it
	ctx := log.ContextWithLogger(context.Background(), logger)

	errChan := make(chan error, 1)

	err := app.setupFunc(ctx, app)
	if err != nil {
		logger.WithError(err).Error("failed to setup app")
		os.Exit(1)
	}

	if len(app.apis) == 0 {
		logger.Error("the app can only work having http apis for now")
		os.Exit(1)
	}

	// starts the HTTP server that offers health-checking
	go app.startHealthServer(ctx, errChan)

	// starts the HTTP server to serve API requests
	go app.startAPIServer(ctx, errChan)

	// waits for a shutdown signal or an error in the current goroutine
	app.waitForShutdown(ctx, errChan)
}

func (app *StandardApp) WithHTTPAPI(api http.API) {
	api.RegisterRoutes(app.router)

	app.apis = append(app.apis, api)
}

func (app *StandardApp) WithPostgresConnection(db string) (*sql.DB, error) {
	source := "sslmode=disable timezone=UTC user=postgres password=test dbname=" + db

	conn, err := sql.Open("postgres", source)
	if err != nil {
		return nil, fmt.Errorf("error connecting with PGUSER, PGPASSWORD to %q: %v", db, err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping db: %v", err)
	}

	return conn, nil
}

// startHealthServer starts a Health providing server for e.g. kubernetes liveness/readiness probe
func (app *StandardApp) startHealthServer(ctx context.Context, errChan chan error) {

	router := basicRouter()
	router.HandleFunc("/healthcheck", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		w.WriteHeader(gohttp.StatusOK)
	})

	addr := ":8585"
	httpServer := newHTTPServer(router, addr)
	app.toShutdown = append(app.toShutdown, httpServer)

	log.FromContext(ctx).
		WithField("server", "health").
		WithField("addr", addr).Info("starting server")

	errChan <- httpServer.Start(ctx)
}

// startAPIServer starts the server to serve the registered API
func (app *StandardApp) startAPIServer(ctx context.Context, errChan chan error) {

	addr := ":8080"
	httpServer := newHTTPServer(app.router, addr)
	app.toShutdown = append(app.toShutdown, httpServer)

	log.FromContext(ctx).
		WithField("server", "api").
		WithField("addr", addr).Info("starting server")

	errChan <- httpServer.Start(ctx)
}

// waitForShutdown blocks the current goroutine until a stop signal is received or a Server returns an error, also
// handling the shutdown of the dependant components (http servers etc)
func (app *StandardApp) waitForShutdown(ctx context.Context, errChan chan error) {
	// set up signal handlers
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(signals)

	logger := log.FromContext(ctx)

	select {
	case err := <-errChan:
		logger.WithError(err).Error("received error from one of the components")
		os.Exit(1)
	case sig := <-signals:
		logger.WithField("sig", sig).Info("signal received, stopping gracefully")
	}

	for _, shutdown := range app.toShutdown {
		if err := shutdown.Shutdown(ctx); err != nil {
			logger.WithError(err).Error("error stopping server")
		}
	}

	os.Exit(0)
}

func basicRouter() *mux.Router {
	// a custom Router that traces requests could be added here for monitoring/instrumentation
	return mux.NewRouter()
}
