package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aditya-dl/QuickShare/backend/api"
	"github.com/aditya-dl/QuickShare/backend/store"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// TODO: add env variables here
	listenAddr := ":8080"
	uploadDir := "./uploads" // relative to where the binary runs

	// initialize the store
	dataStore, err := store.NewMemoryStore(uploadDir)
	if err != nil {
		log.Fatalf("Failed to initialize memory store: %v", err)
	}

	// initialize the api handler with the store
	apiHandler := &api.API{Store: dataStore}

	// setup router
	r := mux.NewRouter()

	// API routes under /api prefix
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/snippets", apiHandler.CreateTextSnippetHandler).Methods("POST")
	apiRouter.HandleFunc("/files", apiHandler.UploadFileHandler).Methods("POST")
	apiRouter.HandleFunc("/files/{id}/download", apiHandler.DownloadFileHandler).Methods("GET")
	apiRouter.HandleFunc("/items", apiHandler.ListItemsHandler).Methods("GET")
	apiRouter.HandleFunc("/items/{id}", apiHandler.GetItemHandler).Methods("GET")
	apiRouter.HandleFunc("/items/{id}", apiHandler.DeleteItemHandler).Methods("DELETE")

	// --- CORS handling (essential for frontend dev server) ---
	// TODO: Adjust origins, methods, headers as needed for production
	// Allowing requests from default Next.js dev server port 3000
	corsOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"}) // TODO: add deployed frontend url later
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}) // TODO: adjust if using auth

	// Apply CORS middleware to the main router
	corsRouter := handlers.CORS(corsOrigins, corsMethods, corsHeaders)(r)

	// Setup server
	srv := &http.Server{
		Handler: corsRouter,
		Addr: listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Printf("Server starting on %s", listenAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}