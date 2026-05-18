package photo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/strbagus/gof-gallery/internal/database"
	req "github.com/strbagus/gof-gallery/pkg/request"
	res "github.com/strbagus/gof-gallery/pkg/response"
)

func GetListPhoto(ctx context.Context, param *req.Pagination) ([]Photo, res.Metadata, error) {
	result := make([]Photo, 0)
	var meta res.Metadata

	whereClause := ""
	args := []any{param.Limit}
	argCount := 1

	if param.Search != "" {
		argCount++
		whereClause += fmt.Sprintf(" and (p.filename ilike $%d or e.name ilike $%d)", argCount, argCount)
		args = append(args, "%"+param.Search+"%")
	}

	query := fmt.Sprintf(`
		select
			p.id,
			p.event_id,
			e.slug as event_slug,
			p.filename,
			p.original_url,
			p.previews,
			count(*) over()
		from
			public.photos p
		join
			public.events e on e.id = p.event_id
		%s
		order by p.%v %v
		limit $1 offset $%d
	`, whereClause, param.OrderBy, param.OrderDir, argCount+1)

	offset := (param.Page - 1) * param.Limit
	args = append(args, offset)

	rows, err := database.PgxPool.Query(ctx, query, args...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()
	for rows.Next() {
		var item Photo
		var previewsJSON []byte
		if err := rows.Scan(&item.ID, &item.EventID, &item.EventSlug, &item.Filename, &item.OriginalURL, &previewsJSON); err != nil {
			return nil, meta, err
		}
		if err := json.Unmarshal(previewsJSON, &item.Previews); err != nil {
			item.Previews = PhotoPreviews{}
		}
		result = append(result, item)
	}

	totalQuery := `
		select
			count(*)
		from
			public.photos p
		limit 1
	`
	var totalRecords int

	err = database.PgxPool.QueryRow(ctx, totalQuery).
		Scan(&totalRecords)

	if err != nil {
		return nil, meta, err
	}

	meta = res.Metadata{
		Total:    totalRecords,
		Page:     param.Page,
		Limit:    param.Limit,
		OrderBy:  param.OrderBy,
		OrderDir: param.OrderDir,
	}

	return result, meta, nil
}

func CreatePhoto(ctx context.Context, req *CreatePhotoReq) error {
	var eventID int
	err := database.PgxPool.QueryRow(ctx, "select id from public.events where slug = $1", req.EventSlug).Scan(&eventID)
	if err != nil {
		return fmt.Errorf("event with slug %s not found: %w", req.EventSlug, err)
	}

	previewsJSON, err := json.Marshal(req.Previews)
	if err != nil {
		return err
	}

	_, err = database.PgxPool.Exec(ctx, `
		insert into public.photos (event_id, event_slug, filename, original_url, previews, created_at)
		values ($1, $2, $3, $4, $5, now())
	`, eventID, req.EventSlug, req.Filename, req.OriginalURL, previewsJSON)

	return err
}
