package main

import (
	"errors"
	"net"
	"net/http"
	"os"
	"time"
	"log"


	"github.com/caring/articles/internal/db"

	"github.com/caring/articles/pb"
	"github.com/caring/go-packages/pkg/logging"
	"github.com/caring/go-packages/pkg/tracing"
	"github.com/getsentry/sentry-go"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)


var (
	store        *db.Store
	dbConnection string
)


var (
	l *logging.Logger
	t *tracing.Tracer
)

var (
	g *grpc.Server
)

func init() {
	l = initLogger()
	initSentry(l)

	dbConnection = setDBConnectionString(l)
	migrateDatabase(l, dbConnection)
	store = initStore(l, dbConnection)

	t = initTracing(l)
	g = createGRPCServer(l, t)
}

func main() {
	defer sentry.Flush(5 * time.Second)
	defer t.Close()
	defer l.Sync()
	defer l.Close()

	// main listener
	lis, err := net.Listen("tcp", ":"+envMust("PORT"))
	if err != nil {
		sentry.CaptureException(err)
		l.Fatal("Failed to initialize net listener:" + err.Error())
	}

	// create a cmux
	m := cmux.New(lis)
	// match connections in order:
	// first grpc, then http.
	grpcL := m.Match(cmux.HTTP2())
	httpL := m.Match(cmux.HTTP1Fast())

	// register the server with gRPC
	pb.RegisterArticlesServer(g, &service{})

	// Add a health check endpoint for automated container monitoring
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// make an error channel to collect the exits of each protocol's Serve()
	eChan := make(chan error)

	// start listeners for each protocol
	go func() { eChan <- g.Serve(grpcL) }()
	go func() { eChan <- http.Serve(httpL, nil) }()

	// all systems are a go
	l.Info("server started: multiplexed http/1, http/2",
		logging.String("port", envMust("PORT")),
		logging.String("multiplexed", "true"),
	)

	// serve it up
	go func() { eChan <- m.Serve() }()

	for err := range eChan {
		if err != nil {
			sentry.CaptureException(err)
			l.Error("Error from one of the HTTP protocols:" + err.Error())
		}
	}

}

// fetches and returns the given env variable, fatals and
// captures an exception if the variable is an empty string
func envMust(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		e := errors.New("environment variable missing - " + varName)
		sentry.CaptureException(e)
		if l != nil {
			l.Fatal(e.Error())
		} else {
			log.Fatalln(e.Error())
		}
	}
	return value
}
