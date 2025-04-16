package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"outerspace-go/lib/grpc"
)

func main() {
	// Get server address from environment variable or use default
	serverAddr := os.Getenv("GRPC_SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "mattintosh.local:50053"
	}

	// Create a new client
	client, err := grpc.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get latest launch
	fmt.Println("\n=== Getting Latest Launch ===")
	launch, err := client.GetLatestLaunch(ctx)
	if err != nil {
		log.Printf("Failed to get latest launch: %v", err)
	} else {
		fmt.Printf("Flight Number: %d\n", launch.FlightNumber)
		fmt.Printf("Mission Name: %s\n", launch.MissionName)
		fmt.Printf("Date (UTC): %s\n", launch.DateUtc)
		fmt.Printf("Success: %v\n", launch.Success)
		fmt.Printf("Details: %s\n", launch.Details)
	}

	// Get all rockets
	fmt.Println("\n=== Getting All Rockets ===")
	rockets, err := client.GetRockets(ctx)
	if err != nil {
		log.Printf("Failed to get rockets: %v", err)
	} else {
		fmt.Println("Available Rockets:")
		for _, rocket := range rockets.Rockets {
			fmt.Printf("- %s (ID: %s)\n", rocket.Name, rocket.Id)
		}
	}

	// Get specific rocket details (using the first rocket ID from the list)
	if rockets != nil && len(rockets.Rockets) > 0 {
		fmt.Println("\n=== Getting Rocket Details ===")
		rocketID := rockets.Rockets[0].Id
		rocket, err := client.GetRocket(ctx, rocketID)
		if err != nil {
			log.Printf("Failed to get rocket details: %v", err)
		} else {
			fmt.Printf("Rocket Name: %s\n", rocket.Name)
			fmt.Printf("Description: %s\n", rocket.Description)
			fmt.Printf("Height: %.2f meters\n", rocket.HeightMeters)
			fmt.Printf("Mass: %d kg\n", rocket.MassKg)
		}
	}

	// Get math fact
	fmt.Println("\n=== Getting Math Fact ===")
	mathFact, err := client.GetMathFact(ctx)
	if err != nil {
		log.Printf("Failed to get math fact: %v", err)
	} else {
		fmt.Printf("Number: %d\n", mathFact.Number)
		fmt.Printf("Type: %s\n", mathFact.Type)
		fmt.Printf("Fact: %s\n", mathFact.Text)
		fmt.Printf("Found: %v\n", mathFact.Found)
	}
}
