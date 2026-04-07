package destination

type PayloadDestination struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Gmaps       string  `json:"gmaps"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type ResDestination struct {
	ID string `json:"id"`
	PayloadDestination
	CreatedAt    string `json:"created_at"`
	LastModified string `json:"last_modified"`
}
