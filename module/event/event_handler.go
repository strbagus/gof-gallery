package event

import (
	"errors"
	req "github.com/strbagus/gof-gallery/pkg/request"
	res "github.com/strbagus/gof-gallery/pkg/response"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
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
