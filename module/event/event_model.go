package event

import (
	"time"
)

type Event struct {
	Name        string    `json:"name" db:"name" validate:"required,min=3"`
	Slug        string    `json:"slug" db:"slug" validate:"required"`
	Location    *string   `json:"location" db:"location" validate:"omitempty"`
	Date        *time.Time `json:"date" db:"date" validate:"omitempty"`
	Description *string   `json:"description" db:"description" validate:"omitempty"`
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
