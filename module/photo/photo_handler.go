package photo

import (
	"errors"
	"slices"
	// "strings"

	// "github.com/strbagus/gof-gallery/module/event"
	req "github.com/strbagus/gof-gallery/pkg/request"
	res "github.com/strbagus/gof-gallery/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	// "github.com/golang-jwt/jwt/v5"
)

// ListPhoto handles listing photos for a specific event slug.
// @Summary List photos by event slug
// @Description Get a paginated list of photos associated with an event slug. If the event is private, a valid access key (JWT) must be provided in the Authorization header.
// @Tags photos
// @Accept json
// @Produce json
// @Param slug path string true "Event Slug"
// @Param Authorization header string false "JWT access key (Bearer <token>)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} response.JSONResponse{data=[]photo.Photo,metadata=response.Metadata} "Successfully retrieved photos"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 401 {object} response.JSONResponse "Unauthorized"
// @Failure 404 {object} response.JSONResponse "Event not found"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Router /photos/{slug} [get]
func ListPhoto(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return res.Error(c, fiber.StatusBadRequest, "Slug event tidak boleh kosong", nil)
	}

	/* // Check if event is private and get salt
	salt, isPrivate, err := event.GetEventSaltBySlug(c.Context(), slug)
	if err != nil {
		return res.Error(c, fiber.StatusNotFound, "Event tidak ditemukan", err.Error())
	}

	if isPrivate {
		authHeader := c.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" {
			return res.Error(c, fiber.StatusUnauthorized, "Event ini privat, silakan masukkan access key", nil)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signing method")
			}
			return []byte(salt), nil
		})

		if err != nil || !token.Valid {
			return res.Error(c, fiber.StatusUnauthorized, "Access key tidak valid", nil)
		}

		// Optionally check if the token belongs to this slug
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["slug"] != slug {
				return res.Error(c, fiber.StatusUnauthorized, "Access key tidak valid untuk event ini", nil)
			}
		} else {
			return res.Error(c, fiber.StatusUnauthorized, "Claims tidak valid", nil)
		}
	} */

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

	list, meta, err := GetListPhoto(c.Context(), slug, f)

	if err != nil {
		return res.Error(c, fiber.StatusInternalServerError, "Terjadi kesalahan saat melakukan query", err.Error())
	}

	return res.Success(c, "Berhasil mendapatkan data", list, &meta)
}

// AddPhoto adds a new photo to an event.
// @Summary Add a new photo
// @Description Add a new photo with its original URL and previews to an event by its slug.
// @Tags photos
// @Accept json
// @Produce json
// @Param request body photo.CreatePhotoReq true "Photo data"
// @Success 200 {object} response.JSONResponse "Successfully added photo"
// @Failure 400 {object} response.JSONResponse "Bad request"
// @Failure 500 {object} response.JSONResponse "Internal server error"
// @Router /photos [post]
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
