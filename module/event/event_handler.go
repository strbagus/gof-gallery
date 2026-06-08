package event

import (
	"errors"
	req "github.com/strbagus/gof-gallery/pkg/request"
	res "github.com/strbagus/gof-gallery/pkg/response"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// ListEvent handles the listing of events with pagination and search.
// @Summary List all events
// @Description Get a paginated list of events with optional search and sorting.
// @Tags events
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param search query string false "Search by name or slug"
// @Param order_by query string false "Column to order by" default(created_at)
// @Param order_dir query string false "Order direction (asc/desc)" default(desc)
// @Success 200 {object} response.JSONResponse{data=[]event.EventResponse,metadata=response.Metadata} "Successfully retrieved events"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Router /events [get]
func ListEvent(c fiber.Ctx) error {
	f := new(req.Pagination)

	if err := c.Bind().Query(f); err != nil {
		return res.Error(c, fiber.StatusBadRequest, "Format parameter salah", err.Error())
	}
	wlColumn := [...]string{"id", "name", "slug", "created_at"}
	f.SetDefaults("created_at", "desc")

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(f); err != nil {
		return res.Error(c, fiber.ErrBadRequest.Code, "Validasi parameter gagal", err.Error())
	}
	if !slices.Contains(wlColumn[:], f.OrderBy) {
		return res.Error(c, fiber.ErrBadRequest.Code, "Validasi parameter gagal", errors.New("Error:Field validation for 'OrderBy'. can't use column '"+f.OrderBy+"' as order by").Error())
	}

	list, meta, err := GetListEvent(c.Context(), f)

	if err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Terjadi kesalahan saat melakukan query", err.Error())
	}

	return res.Success(c, "Berhasil mendapatkan data", list, &meta)
}

// AddEvent handles the creation of a new event.
// @Summary Create a new event
// @Description Create a new event with the provided details.
// @Tags events
// @Accept json
// @Produce json
// @Param request body event.CreateRequest true "Event data"
// @Success 200 {object} response.JSONResponse "Successfully created event"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Router /events [post]
func AddEvent(c fiber.Ctx) error {
	reqBody := new(CreateRequest)

	if err := c.Bind().JSON(reqBody); err != nil {
		return res.Error(c, fiber.StatusBadRequest, "Format body salah", err.Error())
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(reqBody); err != nil {
		return res.Error(c, fiber.ErrBadRequest.Code, "Validasi data gagal", err.Error())
	}

	if err := CreateEvent(c.Context(), reqBody); err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Gagal membuat event", err.Error())
	}

	return res.Success(c, "Berhasil membuat event", nil, nil)
}

// GenerateAccessKey generates a JWT access key for a private event using its salt.
// @Summary Generate access key for private event
// @Description Authenticate access to a private event and return a JWT access key signed with the event's salt.
// @Tags events
// @Accept json
// @Produce json
// @Param slug path string true "Event Slug"
// @Param request body event.AccessKeyRequest true "Password"
// @Success 200 {object} response.JSONResponse{data=event.AccessKeyResponse} "Successfully generated access key"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 401 {object} response.JSONResponse "Unauthorized"
// @Failure 404 {object} response.JSONResponse "Event not found"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Router /events/{slug}/access-key [post]
func GenerateAccessKey(c fiber.Ctx) error {
	slug := c.Params("slug")
	reqBody := new(AccessKeyRequest)

	if err := c.Bind().JSON(reqBody); err != nil {
		return res.Error(c, fiber.StatusBadRequest, "Format body salah", err.Error())
	}

	salt, isPrivate, err := GetEventSaltBySlug(c.Context(), slug)
	if err != nil {
		return res.Error(c, fiber.StatusNotFound, "Event tidak ditemukan", err.Error())
	}

	if !isPrivate {
		return res.Error(c, fiber.StatusBadRequest, "Event ini bukan event privat", nil)
	}

	// For this task, we sign the JWT with the event's unique salt.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"slug": slug,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(salt))
	if err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Gagal generate access key", err.Error())
	}

	return res.Success(c, "Berhasil generate access key", AccessKeyResponse{AccessKey: tokenString}, nil)
}
