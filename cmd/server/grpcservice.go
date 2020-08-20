package main

import (
  _ "context"

  _ "github.com/caring/articles/internal/db"

  _ "github.com/caring/articles/internal/handlers"
  "github.com/caring/articles/pb"
  _ "github.com/caring/go-packages/pkg/errors"
  _ "google.golang.org/grpc/codes"
)

type service struct {
}

func (s *service) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
    l.Printf("Received: %v", in.Data)
    resp := "Data: " + in.Data
    
    ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
    defer cancel()
    status := "up"
    if err := store.Ping(ctx); err != nil {
      status = "down"
    }
    return &pb.PingResponse{Data: resp + "; Database: " + status}, nil
    
}