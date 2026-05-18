package event

import (
	"context"
	"fmt"
	"github.com/strbagus/gof-gallery/internal/database"
	req "github.com/strbagus/gof-gallery/pkg/request"
	res "github.com/strbagus/gof-gallery/pkg/response"
)

func GetListEvent(ctx context.Context, param *req.Pagination) ([]Event, res.Metadata, error) {
	result := make([]Event, 0)
	var meta res.Metadata
	
	whereClause := ""
	args := []any{param.Limit}
	argCount := 1
	
	if param.Search != "" {
		argCount++
		whereClause += fmt.Sprintf(" and (e.name ilike $%d or e.slug ilike $%d)", argCount, argCount)
		args = append(args, "%"+param.Search+"%")
	}

	query := fmt.Sprintf(`
		select
			e.id,
			e."name",
			e.slug
		from
			public.events e
		%s
		order by e.%v %v
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
		var item Event
		if err := rows.Scan(&item.ID, &item.Name, &item.Slug); err != nil {
			return nil, meta, err
		}
		result = append(result, item)
	}

	totalQuery := `
		select
			count(*)
		from
			public.events e
		limit 1
	`
	var totalRecords int

	err = database.PgxPool.QueryRow(ctx, totalQuery).
		Scan(&totalRecords)

	if err != nil {
		return nil, meta, err
	}

	meta = res.Metadata{
		Total:     totalRecords,
		Page:      param.Page,
		Limit:     param.Limit,
		OrderBy:   param.OrderBy,
		OrderDir:  param.OrderDir,
	}

	return result, meta, nil
}
