package photo

import (
	"context"
	"fmt"
	"github.com/strbagus/gof-gallery/internal/database"
	req "github.com/strbagus/gof-gallery/pkg/request"
	res "github.com/strbagus/gof-gallery/pkg/response"
)

func GetListPhoto(ctx context.Context, slug string, param *req.Pagination) ([]string, res.Metadata, error) {
	result := make([]string, 0)
	var meta res.Metadata

	whereClause := "where p.event_slug = $1"
	filterArgs := []any{slug}

	limitIdx := len(filterArgs) + 1
	offsetIdx := len(filterArgs) + 2

	query := fmt.Sprintf(`
		select
			p.preview
		from
			public.photos p
		%s
		limit $%d offset $%d
	`, whereClause, limitIdx, offsetIdx)

	mainArgs := append(filterArgs, param.Limit, (param.Page-1)*param.Limit)

	rows, err := database.PgxPool.Query(ctx, query, mainArgs...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()
	for rows.Next() {
		// var item Photo
		var filename string
		if err := rows.Scan(
			&filename,
		); err != nil {
			return nil, meta, err
		}
		result = append(result, filename)
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
		Total: totalRecords,
		Page:  param.Page,
		Limit: param.Limit,
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
		insert into public.photos (event_id, event_slug, filename, preview, created_at)
		values ($1, $2, $3, $4, now())
	`, eventID, req.EventSlug, req.Filename, req.Preview)

	return err
}

// GetPhotoBySlugAndPreview retrieves a photo by event slug and preview filename
func GetPhotoBySlugAndPreview(ctx context.Context, slug string, preview string) (*Photo, error) {
	var item Photo
	query := `
		select 
			p.event_slug, 
			p.filename 
		from 
			public.photos p 
		where 
			p.event_slug = $1 
			and p.preview = $2 
		limit 1
	`
	err := database.PgxPool.QueryRow(ctx, query, slug, preview).Scan(
		&item.EventSlug,
		&item.Filename,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

