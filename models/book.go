package models

import "time"

type Book struct {
	ID        *string   `json:"id,omitempty" bson:"_id,omitempty"`
	Title     *string   `json:"title"`
	Author    *string   `json:"author"`
	Year      *int      `json:"year"`
	Publisher *string   `json:"publisher"`
	PageCount *int      `json:"pageCount"`
	IsReading *bool     `json:"isReading"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
