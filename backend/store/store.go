package store

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/aditya-dl/QuickShare/backend/models"
	"github.com/google/uuid"
)

// Store deifnes the interface for storing shared items
type Store interface {
	AddItem(item models.SharedItem, fileData multipart.File) (models.SharedItem, error)
	GetItem(id string) (models.SharedItem, bool)
	ListItems() []models.SharedItem
	DeleteItem(id string) bool
	GetFilePath(id string) (string, string, bool) // Returns file path, original filename, found status
}

type MemoryStore struct {
	mu        sync.RWMutex                 // Mutex for thread-safe access to items map
	items     map[string]models.SharedItem // in-memory storage for item metadata
	UploadDir string                       // directory where uploaded files are stored
}

// NewMemoryStore initializes a new MemoryStore with the specified upload directory
func NewMemoryStore(uploadDir string) (*MemoryStore, error) {
	// Ensure upload directory is absolute and exists
	absUploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for upload directory: %w", err)
	}
	if err := os.MkdirAll(absUploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}
	log.Printf("Using upload directory: %s", absUploadDir)

	return &MemoryStore{
		items:     make(map[string]models.SharedItem),
		UploadDir: absUploadDir,
	}, nil
}

// generateNameFromContent generates a name for the snippet based on its content
func generateNameFromContent(content string) (name string) {
	if content == "" {
		return "Untitled Snippet"
	}
	words := strings.Fields(content)
	numWords := 0

	for _, word := range words {
		if utf8.RuneCountInString(name)+utf8.RuneCountInString(word)+1 > 50 && numWords > 0 {
			break
		}
		if name != "" {
			name += " "
		}
		name += word
		numWords++
		if numWords >= 7 {
			break
		}
	}
	if len(name) > 0 && len(name) < len(content) && len(strings.TrimSpace(content)) > len(name) {
		return name + "..."
	}
	return name
}

// AddItem adds a new item to the store
func (s *MemoryStore) AddItem(item models.SharedItem, fileData multipart.File) (models.SharedItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item.ID = uuid.NewString()
	item.CreatedAt = time.Now()

	// Handle automatic naming
	if item.Type == models.ItemTypeText && item.Name == "" && item.Content != "" {
		item.Name = generateNameFromContent(item.Content)
	} else if item.Type == models.ItemTypeFile && item.Name == "" {
		item.Name = filepath.Base(item.FileName)
	}

	// Handle file uploads
	if item.Type == models.ItemTypeFile {
		if fileData == nil {
			return models.SharedItem{}, fmt.Errorf("file data is required for file items")
		}
		// Generate a unique filename for storage on disk using the item's ID
		// This avoid filename collisions and sanitization issues. 
		// We keep the original filename in item.FileName for download prompts. 
		storageFileName := item.ID + item.Name // e.g., <uuid>.txt
		item.FilePath = filepath.Join(s.UploadDir, storageFileName)

		// Create and write the file
		dst, err := os.Create(item.FilePath)
		if err != nil {
			return models.SharedItem{}, fmt.Errorf("failed to create file: %w", err)
		}
		defer dst.Close()

		size, err := io.Copy(dst, fileData)
		if err != nil {
			// Attempt to clean up the partially written file
			os.Remove(item.FilePath)
			return models.SharedItem{}, fmt.Errorf("failed to write file: %w", err)
		}
		item.Size = size
	}

	s.items[item.ID] = item
	log.Printf("Added item: %s, Type: %s, Name=%s", item.ID, item.Type, item.Name)
	return item, nil
}

// GetItem retrieves an item by its ID
func (s *MemoryStore) GetItem(id string) (models.SharedItem, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, found := s.items[id]
	return item, found
}

// ListItems returns a slice of all current items, sorted by the creation date descending. 
func (s *MemoryStore) ListItems() []models.SharedItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	itemList := make([]models.SharedItem, 0, len(s.items))
	for _, item := range s.items {
		// Create a copy to avoid esposing internal file path in general lists if needed, 
		// though current model omits it via json:"-" anyway. Clear content for files in list view. 
		listItem := item
		if listItem.Type == models.ItemTypeFile {
			listItem.Content = "" // Don't send potentially large text content for files in list
		}
		itemList = append(itemList, listItem)
	}

	// Sort by CreatedAt descending (newest first)
	sort.Slice(itemList , func(i, j int) bool {
		return itemList[i].CreatedAt.After(itemList[j].CreatedAt)
	})

	return itemList
}

// DeleteItem removes an item by its ID and deletes the file if it was a file item
func (s *MemoryStore) DeleteItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, found := s.items[id]
	if !found {
		return fmt.Errorf("item with IDD '%s' not found", id)
	}

	// if it's a file, attempt to delete it from disk first
	if item.Type == models.ItemTypeFile && item.FilePath != "" {
		err := os.Remove(item.FilePath)
		// log error if deletion fails, but proceed to remove metdata anyway
		if err != nil && !os.IsNotExist(err) {
			log.Printf("Warning: Failed to delete file '%s' for item ID '%s': %v", item.FilePath, id, err)
		} else if err == nil {
			log.Printf("Deleted file '%s' for item ID '%s'", item.FilePath, id)
		}
	}

	// delete the item from the map
	delete(s.items, id)
	log.Printf("Deleted item with ID '%s'", id)
	return nil
}

// GetFilePath returns the file path and original filename for a file item
func (s *MemoryStore) GetFilePath(id string) (filePath string, originalName string, found bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, found := s.items[id]
	if !found || item.Type != models.ItemTypeFile {
		return "", "", false
	}
	return item.FilePath, item.FileName, true
}