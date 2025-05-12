package models

import "time"

// ItemType defines whether the item is a text snippet or a file
type ItemType string

const (
	ItemTypeText ItemType = "text" // Represent a text snippet
	ItemTypeFile ItemType = "file" // Represent a file
)

type SharedItem struct {
	ID			string    	`json:"id"`
	Name 		string    	`json:"name"`
	Type 		ItemType  	`json:"type"`
	CreatedAt 	time.Time 	`json:"createdAt"`
	Content 	string    	`json:"content,omitempty"` // For text snippets
	FileName 	string    	`json:"fileName,omitempty"` // For file items
	FilePath 	string    	`json:"-"` // Internal path to stored file, not sent in JSON
	ContentType string    	`json:"contentType,omitempty"` // For file items
	Size 		int64     	`json:"size,omitempty"` // For file items
}