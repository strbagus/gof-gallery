package photo

type PhotoPreviews struct {
	Large  string `json:"lg"`
	Medium string `json:"md"`
	Small  string `json:"sm"`
}

type Photo struct {
	ID          int           `json:"id"`
	EventID     int           `json:"event_id"`
	EventSlug   string        `json:"event_slug"`
	Filename    string        `json:"filename"`
	OriginalURL string        `json:"original_url"`
	Previews    PhotoPreviews `json:"previews"`
}

type CreatePhotoReq struct {
	EventSlug   string        `json:"event_slug" validate:"required"`
	Filename    string        `json:"filename" validate:"required"`
	OriginalURL string        `json:"original_url" validate:"required"`
	Previews    PhotoPreviews `json:"previews" validate:"required"`
}
