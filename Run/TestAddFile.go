package main

import (
	"fmt"
	testframework "github.com/PeernetOfficial/Peernet-test-framework"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Main Spawn test Peernet instances for testing
func main() {
	r := mux.NewRouter()

	// Get Config information
	config, err := testframework.ConfigInit()
	if err != nil {
		fmt.Println(err)
	}

	// TODO: extend for future use-case for a embed dashboard tracker
	srv := &http.Server{
		Handler: r,
		Addr:    config.MainServerAddress,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	manager, err := config.RunManager()
	if err != nil {
		fmt.Println(err)
	}

	// Add test file to Node specified
	testframework.AddFilesInNodes(manager)

	// Lister for the main server
	log.Fatal(srv.ListenAndServe())
}
