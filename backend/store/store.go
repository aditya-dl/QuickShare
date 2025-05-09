package store

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/aditya-dl/QuickShare/backend/models"
	"github.com/google/uuid"
)

type MemoryStore struct {
	mu			sync.RWMutex
	items 		map[string]models.SharedItem
	// For file storage, need a directory path
	UploadDir 	string
}

func NewMemoryStore(uploadDir string) *MemoryStore {
	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("Warning: Could not create upload directory %s: %v\n", uploadDir, err)
	}
	return &MemoryStore{
		items:     make(map[string]models.SharedItem),
		UploadDir: uploadDir,
	}
}

func generateNameFromContent(content string) string {
	if content == "" {
		return "Untitled Snippet"
	}
	words := strings.Fields(content)
	numWords := 0
	name := ""

	for _, word := range words {
		if utf8.RuneCountInString(name) + utf8.RuneCountInString(word) + 1 > 50 && numWords > 0 {
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

func (s *MemoryStore) AddTextSnippet(item models.SharedItem) (models.SharedItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item.ID = uuid.NewString()
	item.CreatedAt = time.Now()
	item.Type = models.ItemTypeText

	if item.Name == "" {
		item.Name = generateNameFromContent(item.Content)
	}

	s.items[item.ID] = item
	return item, nil
}

func (s *MemoryStore) GetItem(id string) (models.SharedItem, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, found := s.items[id]
	if found && (item.ExpiresAt.IsZero() || time.Now().Before(item.ExpiresAt)) {
		return item, true
	}

	// If found but expired, delete it
	if found && !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) {
		delete(s.items, id)
		// TODO: delete associated file from disk if item.Type == models.ItemTypeFile
	}
	return models.SharedItem{}, false
}

func (s *MemoryStore) ListItems() []models.SharedItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var currentItems []models.SharedItem
	for _, item := range s.items {
		if item.ExpiresAt.IsZero() || time.Now().Before(item.ExpiresAt) {
			currentItems = append(currentItems, item)
		} else {
			// TODO: delete associated file from disk if item.Type == models.ItemTypeFile
		}	
	}
	// TODO: sort by CreatedAt descending (newest first)
	return currentItems
}

func (s *MemoryStore) DeleteItem(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, found := s.items[id]
	if found {
		delete(s.items, id)
		if item.Type == models.ItemTypeFile && item.FilePath != "" {
			// TODO: os.Remove(item.FilePath)
		}
		return true
	}

	return false
}