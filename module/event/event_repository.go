package event

import (
	"context"
	"fmt"
	"github.com/strbagus/gof-gallery/internal/database"
	req "github.com/strbagus/gof-gallery/pkg/request"
	res "github.com/strbagus/gof-gallery/pkg/response"
)

func GetListEvent(ctx context.Context, param *req.Pagination) ([]EventResponse, res.Metadata, error) {
	result := make([]EventResponse, 0)
	var meta res.Metadata

	whereClause := ""
	filterArgs := []any{}

	if param.Search != "" {
		whereClause += " and (e.name ilike $1 or e.slug ilike $1)"
		filterArgs = append(filterArgs, "%"+param.Search+"%")
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
			e.is_private
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
		OrderBy:  param.OrderBy,
		OrderDir: param.OrderDir,
	}

	return result, meta, nil
}
