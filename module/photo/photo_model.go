package photo

import "time"


type Photo struct {
	ID          int           `json:"id"`
	EventID     int           `json:"event_id"`
	EventSlug   string        `json:"event_slug"`
	Filename    string        `json:"filename"`
	CreatedAt   time.Time     `json:"created_at"`
}

type CreatePhotoReq struct {
	EventSlug   string        `json:"event_slug" validate:"required"`
	Filename    string        `json:"filename" validate:"required"`
}
