package main

import (
	"log"
	"net/http"
	"outerspace-go/lib"
	"outerspace-go/lib/logger"
)

func main() {
	// Initialize logger
	logger.Init()

	spaceClient := lib.NewSpaceXClient()
	numbersClient := lib.NewNumbersClient()

	// Define routes
	http.HandleFunc("/", lib.HandleRoot())
	http.HandleFunc("/api/latest-launch", lib.HandleLatestLaunch(spaceClient))
	http.HandleFunc("/api/rocket", lib.HandleRocket(spaceClient))
	http.HandleFunc("/api/rockets", lib.HandleListRockets(spaceClient))
	http.HandleFunc("/api/numbers", lib.HandleNumbers(numbersClient))

	// Start server
	log.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
