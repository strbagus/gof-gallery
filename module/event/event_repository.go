package event

import (
	"context"
	"fmt"
	"time"

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

func InsertEventAccessToken(ctx context.Context, token string, eventID int, expiresAt time.Time) (*EventAccessToken, error) {
	var item EventAccessToken
	query := `
		insert into public.event_access_tokens (token, event_id, expires_at)
		values ($1, $2, $3)
		returning id, token, event_id, created_at, expires_at
	`
	err := database.PgxPool.QueryRow(ctx, query, token, eventID, expiresAt).Scan(
		&item.ID,
		&item.Token,
		&item.EventID,
		&item.CreatedAt,
		&item.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func GetActiveTokensByEventSlug(ctx context.Context, slug string) ([]EventAccessToken, error) {
	result := make([]EventAccessToken, 0)
	query := `
		select
			eat.id,
			eat.token,
			eat.event_id,
			eat.created_at,
			eat.expires_at
		from
			public.event_access_tokens eat
		join
			public.events e on eat.event_id = e.id
		where
			e.slug = $1
			and e.deleted_at is null
			and eat.expires_at > now()
		order by
			eat.created_at desc
	`
	rows, err := database.PgxPool.Query(ctx, query, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item EventAccessToken
		if err := rows.Scan(
			&item.ID,
			&item.Token,
			&item.EventID,
			&item.CreatedAt,
			&item.ExpiresAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func DeleteEventAccessToken(ctx context.Context, token string) error {
	_, err := database.PgxPool.Exec(ctx, `
		delete from public.event_access_tokens
		where token = $1
	`, token)
	return err
}

func UpdateEventRepo(ctx context.Context, currentSlug string, req *UpdateRequest) error {
	_, err := database.PgxPool.Exec(ctx, `
		update public.events
		set name = $1, slug = $2, location = $3, date = $4, description = $5, is_private = $6
		where slug = $7 and deleted_at is null
	`, req.Name, req.Slug, req.Location, req.Date, req.Description, req.IsPrivate, currentSlug)
	return err
}

func DeleteEventRepo(ctx context.Context, slug string) error {
	_, err := database.PgxPool.Exec(ctx, `
		update public.events
		set deleted_at = now()
		where slug = $1 and deleted_at is null
	`, slug)
	return err
}


