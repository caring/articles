package main

// This file contains helpers that initialize app insight, developer tooling and database set up that might be run on any given app
import (
	"fmt"
	"log"
	"strconv"

	"github.com/caring/go-packages/pkg/grpc_middleware"
	"github.com/caring/go-packages/pkg/logging"
	"github.com/caring/go-packages/pkg/tracing"
	"github.com/getsentry/sentry-go"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"google.golang.org/grpc"
)

// establish logging from env config
func initLogger() *logging.Logger {
	log.Print("Initializing logger")
	l, err := logging.NewLogger(&logging.Config{})
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal("Error initializing logger:", err.Error())
	}
	log.Print("Done")
	return l
}

// configure sentry from env
func initSentry(logger *logging.Logger) {
	logger.Debug("Initializing Sentry")
	val := envMust("SENTRY_DISABLE")
	disabled, err := strconv.ParseBool(val)
	if err != nil {
		logger.Fatal("Error getting SENTRY_DISABLE variable")
	}
	if disabled {
		logger.Debug("Skipping")
		return
	}

	sentryDsn := envMust("SENTRY_DSN")
	env := envMust("SENTRY_ENV")
	err = sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDsn,
		Environment: env,
	})
	if err != nil {
		logger.Fatal("sentry.Init:" + err.Error())
	}
	logger.Debug("Done")

}

// configure tracing form env
func initTracing(logger *logging.Logger) *tracing.Tracer {
	logger.Debug("Initializing Tracing")
	tracer, err := tracing.NewTracer(&tracing.Config{
		Logger: logger,
	})
	if err != nil {
		sentry.CaptureException(err)
		logger.Fatal("Failed to establish tracing:" + err.Error())
	}
	logger.Debug("Done")
	return tracer
}

// create protocol server with chained interceptors
func createGRPCServer(logger *logging.Logger, tracer *tracing.Tracer) *grpc.Server {
	return grpc.NewServer(
		grpc_middleware.NewGRPCChainedUnaryInterceptor(grpc_middleware.UnaryOptions{
			Logger: logger,
			Tracer: tracer,
		}),
		grpc_middleware.NewGRPCChainedStreamInterceptor(grpc_middleware.StreamOptions{
			Logger: logger,
			Tracer: tracer,
		}),
	)
}


// create the db connection string from env
func setDBConnectionString(logger *logging.Logger) string {
	logger.Debug("Creating DB connection string")
	user := envMust("DB_USER")
	pwd := envMust("DB_PWD")
	host := envMust("DB_HOST")
	port := envMust("DB_PORT")
	schema := envMust("DB_SCHEMA")
	logger.Debug("Done")
	return user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + schema
}

// perform the database migration from env config
func migrateDatabase(logger *logging.Logger, connectionString string) {
	logger.Info("Connecting to DB")

	m, err := migrate.New(
		envMust("DB_MIGRATIONS_SRC"),
		"mysql://"+connectionString,
	)
	if err != nil {
		sentry.CaptureException(err)
		logger.Fatal("Failure running migrations to update database:" + err.Error())
	}
	logger.Info("Running migration")
	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			sentry.CaptureException(err)
			logger.Fatal("Migrations Failed: " + err.Error())
		}
	}
	version, dirty, mErr := m.Version()
	logger.Info(fmt.Sprint("Current migration version: ", version))
	logger.Info(fmt.Sprint("Migration dirty: ", dirty))
	if mErr != nil {
		sentry.CaptureException(err)
		logger.Fatal("Migration error: " + mErr.Error())
	}
	logger.Debug("Done")
}

