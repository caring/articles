package main

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/caring/articles/pb"
    "google.golang.org/grpc"
)

const (
    defaultAddress = "localhost:"+envMust("PORT")
    defaultData    = "00"
)

func main() {
    data := defaultData
    address := defaultAddress
    if len(os.Args) > 2 {
        address = os.Args[1]
        data = os.Args[2]
    }

    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewArticlesServiceClient(conn)

    index := 0
    for {
        trip_time := time.Now()
        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
        r, err := c.Ping(ctx, &pb.PingRequest{Data: data})
        if err != nil {
            log.Fatalf("could not connect to: %v", err)
        }

        log.Printf("%d characters roundtrip to (%s): seq=%d time=%s", len(data), address, index, time.Since(trip_time))
        log.Print(r.Data)
        time.Sleep(1 * time.Second)
        index++
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