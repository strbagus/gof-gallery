package event

import (
	"context"
	"fmt"
	"github.com/strbagus/gof-gallery/internal/database"
	res "github.com/strbagus/gof-gallery/pkg/response"
)

func GetListEvent(ctx context.Context, param *ListEventRequest) ([]EventResponse, res.Metadata, error) {
	result := make([]EventResponse, 0)
	var meta res.Metadata

	whereClause := ""
	filterArgs := []any{}

	if param.Search != "" {
		whereClause += " and (e.name ilike $1 or e.slug ilike $1)"
		filterArgs = append(filterArgs, "%"+param.Search+"%")
	}

	if param.IsPrivate != nil {
		isPrivateVal := *param.IsPrivate == 1
		placeholderIdx := len(filterArgs) + 1
		whereClause += fmt.Sprintf(" and e.is_private = $%d", placeholderIdx)
		filterArgs = append(filterArgs, isPrivateVal)
	}

	limitIdx := len(filterArgs) + 1
	offsetIdx := len(filterArgs) + 2

	query := fmt.Sprintf(`
		select
			e.id,
			e.created_at,
			e."name",
			e.slug,
			e.location,
			e.date,
			e.description,
			e.is_private,
			(select count(*) from public.photos p where p.event_id = e.id) as total_photos,
			(select p.preview from public.photos p where p.event_id = e.id order by p.created_at asc limit 1) as thumbnail
		from
			public.events e
		where e.deleted_at is null %s
		order by e.%v %v
		limit $%d offset $%d
	`, whereClause, param.OrderBy, param.OrderDir, limitIdx, offsetIdx)

	mainArgs := append(filterArgs, param.Limit, (param.Page-1)*param.Limit)

	rows, err := database.PgxPool.Query(ctx, query, mainArgs...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()
	for rows.Next() {
		var item EventResponse
		if err := rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&item.Name,
			&item.Slug,
			&item.Location,
			&item.Date,
			&item.Description,
			&item.IsPrivate,
			&item.TotalPhotos,
			&item.Thumbnail,
		); err != nil {
			return nil, meta, err
		}
		result = append(result, item)
	}

	totalQuery := fmt.Sprintf(`
		select
			count(*)
		from
			public.events e
		where e.deleted_at is null %s
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
		OrderBy:  &param.OrderBy,
		OrderDir: &param.OrderDir,
	}

	return result, meta, nil
}

func CreateEvent(ctx context.Context, req *CreateRequest) error {
	_, err := database.PgxPool.Exec(ctx, `
		insert into public.events (name, slug, location, date, description, is_private, created_at)
		values ($1, $2, $3, $4, $5, $6, now())
	`, req.Name, req.Slug, req.Location, req.Date, req.Description, req.IsPrivate)

	return err
}

func GetEventSaltBySlug(ctx context.Context, slug string) (string, bool, error) {
	var salt string
	var isPrivate bool
	err := database.PgxPool.QueryRow(ctx, "select salt, is_private from public.events where slug = $1 and deleted_at is null", slug).Scan(&salt, &isPrivate)
	if err != nil {
		return "", false, err
	}
	return salt, isPrivate, nil
}

func GetEventBySlug(ctx context.Context, slug string) (EventDetailResponse, error) {
	var item EventDetailResponse
	query := `
		select
			e.id,
			e.created_at,
			e."name",
			e.slug,
			e.location,
			e.date,
			e.description,
			e.is_private,
			(select count(*) from public.photos p where p.event_id = e.id) as total_photos
		from
			public.events e
		where e.slug = $1 and e.deleted_at is null
	`
	err := database.PgxPool.QueryRow(ctx, query, slug).Scan(
		&item.ID,
		&item.CreatedAt,
		&item.Name,
		&item.Slug,
		&item.Location,
		&item.Date,
		&item.Description,
		&item.IsPrivate,
		&item.TotalPhotos,
	)
	return item, err
}
