package main

import (
	"context"
	"log"
	"net"
	"shippy-service-consignment/proto/consignment"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*consignment.Consignment) (*consignment.Consignment, error)
}

// Repository - Memory repo, simulats a datastore
type Repository struct {
	mu           sync.RWMutex
	consignments []*consignment.Consignment
}

// Create will create a consignment and return it, or return an error
func (repo *Repository) Create(consignment *consignment.Consignment) (*consignment.Consignment, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

// Service should implement all methods to satisfy the service.
type service struct {
	repo repository
}

func (s *service) CreateConsignment(ctx context.Context, req *consignment.Consignment) (*consignment.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &consignment.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	consignment.RegisterShippingServiceServer(s, &service{repo})
	reflection.Register(s)

	log.Println("Running on port: ", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
