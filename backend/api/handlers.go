package api

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aditya-dl/QuickShare/backend/models"
	"github.com/aditya-dl/QuickShare/backend/store"
	"github.com/gorilla/mux"
)

// API provides handlers that depend on the data store.
type API struct {
	Store store.Store
}

// writeJSON encodes data to JSON and writes it to the response writer.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("failed to encode JSON: %v", err)
			http.Error(w, fmt.Sprintf("failed to encode JSON: %v", err), http.StatusInternalServerError)
		}
	}
}

// writeError encodes an error message to JSON and writes it to the response writer.
func writeError(w http.ResponseWriter, status int, message string) {
	log.Printf("API error (%d): %s", status, message)
	writeJSON(w, status, map[string]string{"error": message})
}

// CreateTextSnippetHandler handles POST requests to create a new text snippet.
func (a *API) CreateTextSnippetHandler(w http.ResponseWriter, r *http.Request) {
	var item models.SharedItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if item.Content == "" {
		writeError(w, http.StatusBadRequest, "Content cannot be empty")
		return
	}

	item.Type = models.ItemTypeText
	
	// let the store handle id generation, timestamp, and name generation
	createdItem, err := a.Store.AddItem(item, nil) // no file data for text snippets
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create text snippet")
		return
	}

	writeJSON(w, http.StatusCreated, createdItem)
}

// UploadFileHandler handles POST requests to upload a file.
func (a *API) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// max upload size 
	const maxUploadSize = 100 * 1024 * 1024 // 100 MB
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse multipart form (max size %dMB): %v", maxUploadSize/(1024*1024), err))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid file upload request (missing 'file' field?): "+err.Error())
		return
	}
	defer file.Close()

	// Optional: Get user-provided name from form field, otherwise store uses filename
	itemName := r.FormValue("name")

	item := models.SharedItem {
		Type: models.ItemTypeFile,
		Name: itemName, // can be empty, store will handle default naming
		FileName: handler.Filename,
		ContentType: handler.Header.Get("Content-Type"),
		Size: handler.Size,
	}

	// Let the store handle id generation, timestamp, naming, file saving, and metadata storage
	createdItem, err := a.Store.AddItem(item, file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to upload file: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, createdItem)
}

// ListItemsHandler handles GET requests to list all items.
func (a *API) ListItemsHandler(w http.ResponseWriter, r *http.Request) {
	items := a.Store.ListItems()
	if items == nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve items")
		return
	}

	writeJSON(w, http.StatusOK, items)
}

// GetItemHandler handles GET requests to retrieve a specific item by ID.
func (a *API) GetItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	item, found := a.Store.GetItem(id)
	if !found {
		writeError(w, http.StatusNotFound, "Item not found")
		return
	}
	
	writeJSON(w, http.StatusOK, item)
}

// DownloadFileHandler handles GET requests to download a file by ID.
func (a *API) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	filePath, originalName, found := a.Store.GetFilePath(id)
	if !found {
		writeError(w, http.StatusNotFound, "File not found")
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", originalName))

	// Detect Content-Type if not stored reliably, or use stored type
	contentType := mime.TypeByExtension(filepath.Ext(originalName))
	if contentType == "" {
		contentType = "application/octet-stream" // Fallback to binary
	}
	w.Header().Set("Content-Type", contentType)

	http.ServeFile(w, r, filePath)
	log.Printf("Served file download: %s, Path: %s, Name: %s", id, filePath, originalName)
}

// DeleteItemHandler handles DELETE requests to delete a specific item by ID.
func (a *API) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := a.Store.DeleteItem(id)
	if err != nil {
		// Check if the error indicates "not found" vs. a real delete error
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to delete item: "+err.Error())
		}
		return
	}

	// Success, not content to return
	w.WriteHeader(http.StatusNoContent)
}