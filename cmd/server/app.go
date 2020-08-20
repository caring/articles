package main

// This file contains helpers to initialize application code that is specific to this service
import (
	"github.com/caring/articles/internal/db"
	"github.com/caring/go-packages/pkg/logging"
	"github.com/getsentry/sentry-go"
)


// initialize the store service
func initStore(logger *logging.Logger, connectionString string) *db.Store {
	logger.Debug("Initializing Store")
	// establish a store and connection to the db
	store, err := db.NewStore(connectionString)
	if err != nil {
		sentry.CaptureException(err)
		logger.Fatal("Failed to initialize store:" + err.Error())
	}
	logger.Debug("Store established with database connection")
	return store
}

