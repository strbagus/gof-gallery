package event

import (
	"time"
	req "github.com/strbagus/gof-gallery/pkg/request"
)

type Event struct {
	Name        string     `json:"name" db:"name" validate:"required,min=3"`
	Slug        string     `json:"slug" db:"slug" validate:"required"`
	Location    *string    `json:"location" db:"location" validate:"omitempty"`
	Date        *time.Time `json:"date" db:"date" validate:"omitempty"`
	Description *string    `json:"description" db:"description" validate:"omitempty"`
	IsPrivate   bool       `json:"is_private" db:"is_private"`
}

type EventResponse struct {
	ID          int       `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	TotalPhotos int       `json:"total_photos"`
	Thumbnail   *string   `json:"thumbnail"`
	Event
}

type EventDetailResponse struct {
	ID          int       `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	TotalPhotos int       `json:"total_photos"`
	Event
}

type CreateRequest struct {
	Event
}

type AccessKeyRequest struct {
	Password string `json:"password" validate:"required"`
}

type ListEventRequest struct {
	req.Pagination
	IsPrivate *int `query:"is_private" validate:"omitempty,oneof=0 1"`
}

type EventAccessToken struct {
	ID        int       `json:"id" db:"id"`
	Token     string    `json:"token" db:"token"`
	EventID   int       `json:"event_id" db:"event_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

type GenerateTokenRequest struct {
	DurationDays int `json:"duration_days" validate:"required,min=1"`
}

