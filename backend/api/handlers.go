package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/aditya-dl/QuickShare/backend/models"
	"github.com/aditya-dl/QuickShare/backend/store"
)

type API struct {
	Store *store.MemoryStore
}

func (a * API) CreateTextSnippetHandler(w http.ResponseWriter, r *http.Request) {
	var item models.SharedItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item.Type = models.ItemTypeText
	if item.Name == "" && item.Content != "" {
		words := strings.Fields(item.Content)
		numWords := 0
		name := ""
		for _, word := range words {
			if utf8.RuneCountInString(name) + utf8.RuneCountInString(word) + 1  > 50 && numWords > 0 {
				// Limit length
				break
			}
			if name != "" {
				name += " "
			}
			name += word
			numWords++
			if numWords >= 7 {
				// Limit words
				break
			}
		}
		if len(name) > 0 && len(name) < len(item.Content) { 
			// Add ellipsis if shortened
			item.Name = name + "..."
		} else {
			item.Name = name
		}
	}

	createdItem, err := a.Store.AddItem(item)
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdItem)
}

// ... other handlers ...
// Handle file uploads (multipart forms), saving files, and serving them
// For file uploads, use r.ParseMultipartForm, r.FormFile, io.Copy.
// Store uploaded files in a dedicated directory (e.g., `backend/uploads/`). Ensure this dir is in .gitignore.