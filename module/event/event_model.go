package event

import (
	"time"
)

type Event struct {
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Location    *string   `json:"location" db:"location"`
	Date        *time.Time `json:"date" db:"date"`
	Description *string   `json:"description" db:"description"`
	IsPrivate   bool      `json:"is_private" db:"is_private"`
}

type EventResponse struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Event
}

type CreateRequest struct {
	Event
}
