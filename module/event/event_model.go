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

type AccessKeyResponse struct {
	AccessKey string `json:"access_key"`
}

type ListEventRequest struct {
	req.Pagination
	IsPrivate *int `query:"is_private" validate:"omitempty,oneof=0 1"`
}
