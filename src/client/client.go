package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/relangi/go-grpc/src/proto_out"
	"google.golang.org/grpc"
)

// Connect to the gRPC server
func connect() pb.InventoryClient {
	conn, err := grpc.Dial("localhost:9001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return pb.NewInventoryClient(conn)
}

// Add command
func addItem(sku, name, description string, price float32, quantity uint32) {
	c := connect()

	item := &pb.Item{
		Identifier: &pb.ItemIdentifier{Sku: sku},
		Stock:      &pb.ItemStock{Price: price, Quantity: quantity},
		Information: &pb.ItemInformation{
			Name:        &name,
			Description: &description,
		},
	}

	res, err := c.Add(context.Background(), item)
	if err != nil {
		log.Fatalf("could not add item: %v", err)
	}
	fmt.Printf("Add Response: %s\n", res.Status)
}

// Remove command
func removeItem(sku string) {
	c := connect()

	res, err := c.Remove(context.Background(), &pb.ItemIdentifier{Sku: sku})
	if err != nil {
		log.Fatalf("could not remove item: %v", err)
	}
	fmt.Printf("Remove Response: %s\n", res.Status)
}

// Get command
func getItem(sku string) {
	c := connect()

	item, err := c.Get(context.Background(), &pb.ItemIdentifier{Sku: sku})
	if err != nil {
		log.Fatalf("could not get item: %v", err)
	}
	fmt.Printf("Get Response: %v\n", item)
}

// UpdateQuantity command
func updateQuantity(sku string, change int32) {
	c := connect()

	req := &pb.QuantityChangeRequest{
		Sku:    sku,
		Change: change,
	}

	res, err := c.UpdateQuantity(context.Background(), req)
	if err != nil {
		log.Fatalf("could not update quantity: %v", err)
	}
	fmt.Printf("Update Quantity Response: %v\n", res)
}

// UpdatePrice command
func updatePrice(sku string, price float32) {
	c := connect()

	req := &pb.PriceChangeRequest{
		Sku:   sku,
		Price: price,
	}

	res, err := c.UpdatePrice(context.Background(), req)
	if err != nil {
		log.Fatalf("could not update price: %v", err)
	}
	fmt.Printf("Update Price Response: %v\n", res)
}

// Watch command
func watchItem(sku string) {
	c := connect()

	stream, err := c.Watch(context.Background(), &pb.ItemIdentifier{Sku: sku})
	if err != nil {
		log.Fatalf("could not watch item: %v", err)
	}

	for {
		item, err := stream.Recv()
		if err != nil {
			log.Fatalf("error receiving stream: %v", err)
		}
		fmt.Printf("Received Item: %v\n", item)
		time.Sleep(1 * time.Second) // Adjust as needed
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'add', 'remove', 'get', 'update-quantity', 'update-price', or 'watch' subcommand")
		os.Exit(1)
	}

	subcommand := os.Args[1]
	switch subcommand {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		skuFlag := addCmd.String("sku", "", "SKU of the item")
		nameFlag := addCmd.String("name", "", "Name of the item")
		descFlag := addCmd.String("description", "", "Description of the item")
		priceFlag := addCmd.Float64("price", 0.0, "Price of the item")
		quantityFlag := addCmd.Uint("quantity", 0, "Quantity of the item")
		addCmd.Parse(os.Args[2:])

		if *skuFlag == "" {
			fmt.Println("sku is required")
			os.Exit(1)
		}
		addItem(*skuFlag, *nameFlag, *descFlag, float32(*priceFlag), uint32(*quantityFlag))

	case "remove":
		removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
		skuFlag := removeCmd.String("sku", "", "SKU of the item")
		removeCmd.Parse(os.Args[2:])

		if *skuFlag == "" {
			fmt.Println("sku is required")
			os.Exit(1)
		}
		removeItem(*skuFlag)

	case "get":
		getCmd := flag.NewFlagSet("get", flag.ExitOnError)
		skuFlag := getCmd.String("sku", "", "SKU of the item")
		getCmd.Parse(os.Args[2:])

		if *skuFlag == "" {
			fmt.Println("sku is required")
			os.Exit(1)
		}
		getItem(*skuFlag)

	case "update-quantity":
		updateQuantityCmd := flag.NewFlagSet("update-quantity", flag.ExitOnError)
		skuFlag := updateQuantityCmd.String("sku", "", "SKU of the item")
		changeFlag := updateQuantityCmd.Int("change", 0, "Quantity change")
		updateQuantityCmd.Parse(os.Args[2:])

		if *skuFlag == "" {
			fmt.Println("sku is required")
			os.Exit(1)
		}
		updateQuantity(*skuFlag, int32(*changeFlag))

	case "update-price":
		updatePriceCmd := flag.NewFlagSet("update-price", flag.ExitOnError)
		skuFlag := updatePriceCmd.String("sku", "", "SKU of the item")
		priceFlag := updatePriceCmd.Float64("price", 0.0, "New price")
		updatePriceCmd.Parse(os.Args[2:])

		if *skuFlag == "" {
			fmt.Println("sku is required")
			os.Exit(1)
		}
		updatePrice(*skuFlag, float32(*priceFlag))

	case "watch":
		watchCmd := flag.NewFlagSet("watch", flag.ExitOnError)
		skuFlag := watchCmd.String("sku", "", "SKU of the item")
		watchCmd.Parse(os.Args[2:])

		if *skuFlag == "" {
			fmt.Println("sku is required")
			os.Exit(1)
		}
		watchItem(*skuFlag)

	default:
		fmt.Println("expected 'add', 'remove', 'get', 'update-quantity', 'update-price', or 'watch' subcommand")
		os.Exit(1)
	}
}
