package photo

import (
	"errors"
	"slices"
	req "github.com/strbagus/gof-gallery/pkg/request"
	res "github.com/strbagus/gof-gallery/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func ListPhoto(c fiber.Ctx) error {
	f := new(req.Pagination)

	if err := c.Bind().Query(f); err != nil {
		return res.Error(c, fiber.StatusBadRequest, "Format parameter salah", err.Error())
	}
	wlColumn := [...]string{"id", "event_id", "filename", "created_at"}
	f.SetDefaults("created_at", "desc")

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(f); err != nil {
		return res.Error(c, fiber.ErrBadRequest.Code, "Validasi parameter gagal", err.Error())
	}
	if !slices.Contains(wlColumn[:], f.OrderBy) {
		return res.Error(c, fiber.ErrBadRequest.Code, "Validasi parameter gagal", errors.New("Error:Field validation for 'OrderBy'. can't use column '"+f.OrderBy+"' as order by").Error())
	}

	list, meta, err := GetListPhoto(c.Context(), f)

	if err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Terjadi kesalahan saat melakukan query", err.Error())
	}

	return res.Success(c, "Berhasil mendapatkan data", list, &meta)
}

func AddPhoto(c fiber.Ctx) error {
	reqBody := new(CreatePhotoReq)

	if err := c.Bind().JSON(reqBody); err != nil {
		return res.Error(c, fiber.StatusBadRequest, "Format body salah", err.Error())
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(reqBody); err != nil {
		return res.Error(c, fiber.ErrBadRequest.Code, "Validasi data gagal", err.Error())
	}

	if err := CreatePhoto(c.Context(), reqBody); err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Gagal menambahkan foto", err.Error())
	}

	return res.Success(c, "Berhasil menambahkan foto", nil, nil)
}

