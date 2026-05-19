package photo

import (
	"context"
	"fmt"
	"github.com/strbagus/gof-gallery/internal/database"
	req "github.com/strbagus/gof-gallery/pkg/request"
	res "github.com/strbagus/gof-gallery/pkg/response"
)

func GetListPhoto(ctx context.Context, slug string, param *req.Pagination) ([]Photo, res.Metadata, error) {
	result := make([]Photo, 0)
	var meta res.Metadata

	whereClause := "where p.event_slug = $1"
	filterArgs := []any{slug}

	limitIdx := len(filterArgs) + 1
	offsetIdx := len(filterArgs) + 2

	query := fmt.Sprintf(`
		select
			p.id,
			p.event_id,
			p.event_slug,
			p.filename,
			p.created_at
		from
			public.photos p
		%s
		order by p.%v %v
		limit $%d offset $%d
	`, whereClause, param.OrderBy, param.OrderDir, limitIdx, offsetIdx)

	mainArgs := append(filterArgs, param.Limit, (param.Page-1)*param.Limit)

	rows, err := database.PgxPool.Query(ctx, query, mainArgs...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()
	for rows.Next() {
		var item Photo
		if err := rows.Scan(
			&item.ID,
			&item.EventID,
			&item.EventSlug,
			&item.Filename,
			&item.CreatedAt,
		); err != nil {
			return nil, meta, err
		}
		result = append(result, item)
	}

	totalQuery := fmt.Sprintf(`
		select
			count(*)
		from
			public.photos p
		%s
	`, whereClause)

	var totalRecords int
	err = database.PgxPool.QueryRow(ctx, totalQuery, filterArgs...).Scan(&totalRecords)
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

	_, err = database.PgxPool.Exec(ctx, `
		insert into public.photos (event_id, event_slug, filename, created_at)
		values ($1, $2, $3, $4, $5, now())
	`, eventID, req.EventSlug, req.Filename)

	return err
}
