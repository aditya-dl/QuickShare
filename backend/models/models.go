package models

import "time"

type ItemType string

const (
	ItemTypeText ItemType = "text"
	ItemTypeFile ItemType = "file"
)

type SharedItem struct {
	ID			string    	`json:"id"`
	Name 		string    	`json:"name"`
	Type 		ItemType  	`json:"type"`
	Content 	string    	`json:"content,omitempty"` // For text snippets
	FilePath 	string    	`json:"-"` // Internal path to stored file, not sent in JSON
	FileName 	string    	`json:"fileName,omitempty"` // For file items
	ContentType string    	`json:"contentType,omitempty"` // For file items
	Size 		int64     	`json:"size,omitempty"` // For file items
	CreatedAt 	time.Time 	`json:"createdAt"`
	ExpiresAt 	time.Time 	`json:"expiresAt,omitempty"`
}