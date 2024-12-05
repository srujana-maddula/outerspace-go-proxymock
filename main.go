package main

import (
	"log"
	"net/http"
	"outerspace-go/lib"
)

func main() {
	client := lib.NewSpaceXClient()

	// Define routes
	http.HandleFunc("/", lib.HandleRoot())
	http.HandleFunc("/api/latest-launch", lib.HandleLatestLaunch(client))
	http.HandleFunc("/api/rocket", lib.HandleRocket(client))
	http.HandleFunc("/api/rockets", lib.HandleListRockets(client))

	// Start server
	log.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
