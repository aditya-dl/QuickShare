package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aditya-dl/QuickShare/backend/api"
	"github.com/aditya-dl/QuickShare/backend/store"
	"github.com/gorilla/mux"
)

func main() {
	uploadDir := "./uploads" // Make this cnofigurable (later)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
	absUploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path of upload directory: %v", err)
	}

	dataStore := store.NewMemoryStore(absUploadDir)
	apiHandler := &api.API{Store: dataStore}

	r := mux.NewRouter()

	// API routes
	apiRoutes := r.PathPrefix("/api").Subrouter()
	apiRoutes.HandleFunc("/snippets", apiHandler.CreateTextSnippetHandler).Methods("POST")
	// apiRoutes.HandleFunc("/snippets", apiHandler.ListItemsHandler).Methods("GET") // Adapt to filter by type or have separate endpoints
    // apiRoutes.HandleFunc("/snippets/{id}", apiHandler.GetItemHandler).Methods("GET")
    // apiRoutes.HandleFunc("/snippets/{id}", apiHandler.DeleteItemHandler).Methods("DELETE")

	// apiRoutes.HandleFunc("/files", apiHandler.UploadFileHandler).Methods("POST")
	// apiRoutes.HandleFunc("/files", apiHandler.ListItemsHandler).Methods("GET")
    // apiRoutes.HandleFunc("/files/{id}/download", apiHandler.DownloadFileHandler).Methods("GET") // Specific download route
    // apiRoutes.HandleFunc("/files/{id}", apiHandler.DeleteItemHandler).Methods("DELETE")

	// apiRoutes.HandleFunc("/items", apiHandler.ListItemsHandler).Methods("GET") // General list for both
    // apiRoutes.HandleFunc("/items/{id}", apiHandler.GetItemHandler).Methods("GET")
    // apiRoutes.HandleFunc("/items/{id}", apiHandler.DeleteItemHandler).Methods("DELETE")

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}