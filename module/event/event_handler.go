package event

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	res "github.com/strbagus/gof-gallery/pkg/response"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
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
// @Param is_private query int false "Filter by privacy (0 for public, 1 for private)"
// @Success 200 {object} response.JSONResponse{data=[]event.EventResponse,metadata=response.Metadata} "Successfully retrieved events"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Router /events [get]
func ListEvent(c fiber.Ctx) error {
	f := new(ListEventRequest)

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
// @Security CookieAuth
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

// GetEventDetail handles getting a single event by its slug.
// @Summary Get event by slug
// @Description Get detailed information of a specific event using its slug.
// @Tags events
// @Accept json
// @Produce json
// @Param slug path string true "Event Slug"
// @Success 200 {object} response.JSONResponse{data=event.EventDetailResponse} "Successfully retrieved event"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 404 {object} response.JSONResponse "Event not found"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Security CookieAuth
// @Router /events/{slug} [get]
func GetEventDetail(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return res.Error(c, fiber.StatusBadRequest, "Slug event tidak boleh kosong", nil)
	}

	item, err := GetEventBySlug(c.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res.Error(c, fiber.StatusNotFound, "Event tidak ditemukan", err.Error())
		}
		return res.Error(c, fiber.StatusInternalServerError, "Terjadi kesalahan saat mengambil data event", err.Error())
	}

	return res.Success(c, "Berhasil mendapatkan data event", item, nil)
}

// GenerateToken handles generating a new event access token.
// @Summary Generate event access token
// @Description Generate a secure random access token for a specific event with a given validity duration.
// @Tags event-tokens
// @Accept json
// @Produce json
// @Param slug path string true "Event Slug"
// @Param request body event.GenerateTokenRequest true "Token duration payload"
// @Success 201 {object} response.JSONResponse{data=event.EventAccessToken} "Successfully generated access token"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 404 {object} response.JSONResponse "Event not found"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Security CookieAuth
// @Router /events/{slug}/tokens [post]
func GenerateToken(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return res.Error(c, fiber.StatusBadRequest, "Slug event tidak boleh kosong", nil)
	}

	event, err := GetEventBySlug(c.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res.Error(c, fiber.StatusNotFound, "Event tidak ditemukan", err.Error())
		}
		return res.Error(c, fiber.StatusInternalServerError, "Gagal memverifikasi event", err.Error())
	}

	reqBody := new(GenerateTokenRequest)
	if err := c.Bind().JSON(reqBody); err != nil {
		return res.Error(c, fiber.StatusBadRequest, "Format body salah", err.Error())
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(reqBody); err != nil {
		return res.Error(c, fiber.ErrBadRequest.Code, "Validasi data gagal", err.Error())
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Gagal membuat token acak", err.Error())
	}
	tokenStr := hex.EncodeToString(b)

	expiresAt := time.Now().AddDate(0, 0, reqBody.DurationDays)

	tokenRecord, err := InsertEventAccessToken(c.Context(), tokenStr, event.ID, expiresAt)
	if err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Gagal menyimpan token ke database", err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(res.JSONResponse{
		Status:  "success",
		Message: "Token berhasil dibuat",
		Data:    tokenRecord,
	})
}

// ListTokens handles listing active access tokens for a specific event.
// @Summary List active event tokens
// @Description Retrieve a list of active (non-expired) access tokens for a specific event using its slug.
// @Tags event-tokens
// @Accept json
// @Produce json
// @Param slug path string true "Event Slug"
// @Success 200 {object} response.JSONResponse{data=[]event.EventAccessToken} "Successfully retrieved active tokens"
// @Failure 404 {object} response.JSONResponse "Event not found"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Security CookieAuth
// @Router /events/{slug}/tokens [get]
func ListTokens(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return res.Error(c, fiber.StatusBadRequest, "Slug event tidak boleh kosong", nil)
	}

	_, err := GetEventBySlug(c.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res.Error(c, fiber.StatusNotFound, "Event tidak ditemukan", err.Error())
		}
		return res.Error(c, fiber.StatusInternalServerError, "Gagal memverifikasi event", err.Error())
	}

	tokens, err := GetActiveTokensByEventSlug(c.Context(), slug)
	if err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Gagal mengambil daftar token", err.Error())
	}

	return res.Success(c, "Berhasil mendapatkan daftar token", tokens, nil)
}

// RevokeToken handles deleting/revoking a specific access token.
// @Summary Revoke access token
// @Description Execute a hard delete on a specific access token from the database.
// @Tags event-tokens
// @Accept json
// @Produce json
// @Param slug path string true "Event Slug"
// @Param token path string true "Access Token to revoke"
// @Success 204 "No Content - token successfully deleted"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Security CookieAuth
// @Router /events/{slug}/tokens/{token} [delete]
func RevokeToken(c fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return res.Error(c, fiber.StatusBadRequest, "Token tidak boleh kosong", nil)
	}

	err := DeleteEventAccessToken(c.Context(), token)
	if err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Gagal menghapus token", err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateEvent handles updating an event by its slug.
// @Summary Update event details
// @Description Update the properties of an existing event using its current slug.
// @Tags events
// @Accept json
// @Produce json
// @Param slug path string true "Current Event Slug"
// @Param request body event.UpdateRequest true "Updated event data"
// @Success 200 {object} response.JSONResponse "Successfully updated event"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 404 {object} response.JSONResponse "Event not found"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Security CookieAuth
// @Router /events/{slug} [put]
func UpdateEvent(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return res.Error(c, fiber.StatusBadRequest, "Slug event tidak boleh kosong", nil)
	}

	// Verify event exists
	_, err := GetEventBySlug(c.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res.Error(c, fiber.StatusNotFound, "Event tidak ditemukan", err.Error())
		}
		return res.Error(c, fiber.StatusInternalServerError, "Gagal memverifikasi event", err.Error())
	}

	reqBody := new(UpdateRequest)
	if err := c.Bind().JSON(reqBody); err != nil {
		return res.Error(c, fiber.StatusBadRequest, "Format body salah", err.Error())
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(reqBody); err != nil {
		return res.Error(c, fiber.ErrBadRequest.Code, "Validasi data gagal", err.Error())
	}

	if err := UpdateEventRepo(c.Context(), slug, reqBody); err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Gagal mengupdate event", err.Error())
	}

	return res.Success(c, "Berhasil mengupdate event", nil, nil)
}

// DeleteEvent handles soft-deleting an event.
// @Summary Delete an event
// @Description Soft-delete an event from the system using its slug.
// @Tags events
// @Accept json
// @Produce json
// @Param slug path string true "Event Slug"
// @Success 200 {object} response.JSONResponse "Successfully deleted event"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 404 {object} response.JSONResponse "Event not found"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Security CookieAuth
// @Router /events/{slug} [delete]
func DeleteEvent(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return res.Error(c, fiber.StatusBadRequest, "Slug event tidak boleh kosong", nil)
	}

	// Verify event exists
	_, err := GetEventBySlug(c.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res.Error(c, fiber.StatusNotFound, "Event tidak ditemukan", err.Error())
		}
		return res.Error(c, fiber.StatusInternalServerError, "Gagal memverifikasi event", err.Error())
	}

	if err := DeleteEventRepo(c.Context(), slug); err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Gagal menghapus event", err.Error())
	}

	return res.Success(c, "Berhasil menghapus event", nil, nil)
}


