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
