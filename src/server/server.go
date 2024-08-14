package main

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	pb "github.com/relangi/go-grpc/src/proto_out"
	"google.golang.org/grpc"
)

// server is used to implement the Inventory service
type server struct {
	pb.UnimplementedInventoryServer
	mu      sync.Mutex
	storage map[string]*pb.Item
}

// Add inserts a new Item into the inventory.
func (s *server) Add(ctx context.Context, item *pb.Item) (*pb.InventoryChangeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the item already exists
	if _, exists := s.storage[item.Identifier.Sku]; exists {
		return &pb.InventoryChangeResponse{Status: "Item already exists"}, errors.New("item already exists")
	}

	s.storage[item.Identifier.Sku] = item
	return &pb.InventoryChangeResponse{Status: "Item added"}, nil
}

// Remove removes Items from the inventory.
func (s *server) Remove(ctx context.Context, id *pb.ItemIdentifier) (*pb.InventoryChangeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.storage[id.Sku]; !exists {
		return &pb.InventoryChangeResponse{Status: "Item not found"}, errors.New("item not found")
	}

	delete(s.storage, id.Sku)
	return &pb.InventoryChangeResponse{Status: "Item removed"}, nil
}

// Get retrieves Item information.
func (s *server) Get(ctx context.Context, id *pb.ItemIdentifier) (*pb.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.storage[id.Sku]
	if !exists {
		return nil, errors.New("item not found")
	}

	return item, nil
}

// UpdateQuantity increases or decreases the stock quantity of an Item.
func (s *server) UpdateQuantity(ctx context.Context, req *pb.QuantityChangeRequest) (*pb.InventoryUpdateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.storage[req.Sku]
	if !exists {
		return nil, errors.New("item not found")
	}

	item.Stock.Quantity += uint32(req.Change)

	return &pb.InventoryUpdateResponse{
		Status:   "Quantity updated",
		Price:    item.Stock.Price,
		Quantity: item.Stock.Quantity,
	}, nil
}

// UpdatePrice increases or decreases the price of an Item.
func (s *server) UpdatePrice(ctx context.Context, req *pb.PriceChangeRequest) (*pb.InventoryUpdateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.storage[req.Sku]
	if !exists {
		return nil, errors.New("item not found")
	}

	item.Stock.Price = req.Price

	return &pb.InventoryUpdateResponse{
		Status:   "Price updated",
		Price:    item.Stock.Price,
		Quantity: item.Stock.Quantity,
	}, nil
}

// Watch streams Item updates from the inventory.
func (s *server) Watch(id *pb.ItemIdentifier, stream pb.Inventory_WatchServer) error {
	s.mu.Lock()
	item, exists := s.storage[id.Sku]
	s.mu.Unlock()

	if !exists {
		return errors.New("item not found")
	}

	// Here we just stream the current state of the item every second
	for i := 0; i < 5; i++ {
		err := stream.Send(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInventoryServer(s, &server{storage: make(map[string]*pb.Item)})
	log.Printf("Server is running on port 9001...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
