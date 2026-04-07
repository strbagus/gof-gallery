package destination

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strbagus/travel-agent/internal/database"
	h "strbagus/travel-agent/internal/helper"
	req "strbagus/travel-agent/pkg/request"
	res "strbagus/travel-agent/pkg/response"
)

func GetListDestination(ctx context.Context, param *req.Pagination) ([]ResDestination, res.Metadata, error) {
	result := make([]ResDestination, 0)
	var meta res.Metadata
	query := fmt.Sprintf(`
		select
			d.id,
			d."name",
			coalesce(d.description, ''),
			coalesce(d.gmaps, ''),
			coalesce(d.latitude, 0.00),
			coalesce(d.longitude, 0.00),
			d.created_at,
			coalesce(d.updated_at, d.created_at),
			count(*) over()
		from
			master.destination d
		where
			d.deleted_at is null
		order by d.%v %v
		limit $1 offset $2
	`, param.OrderBy, param.OrderDir)
	// and ($1 = '' or d."name" like '%' || $1 || '%')
	offset := (param.Page - 1) * param.PerPage

	rows, err := database.PgxPool.Query(ctx, query, param.PerPage, offset)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()
	var totalFiltered int
	for rows.Next() {
		var item ResDestination
		var createdAt sql.NullTime
		var lastModified sql.NullTime
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Gmaps, &item.Latitude, &item.Longitude, &createdAt, &lastModified, &totalFiltered); err != nil {
			return nil, meta, err
		}
		item.CreatedAt, err = h.DBToTime(createdAt)
		if err != nil {
			return nil, meta, err
		}
		item.LastModified, err = h.DBToTime(lastModified)
		if err != nil {
			return nil, meta, err
		}
		result = append(result, item)
	}

	totalQuery := `
		select
			count(*)
		from
			master.destination d
		where
			d.deleted_at is null
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
		Filtered:  totalFiltered,
		Page:      param.Page,
		PerPage:   param.PerPage,
		TotalPage: int(math.Ceil(float64(totalFiltered) / float64(param.PerPage))),
		OrderBy:   param.OrderBy,
		OrderDir:  param.OrderDir,
	}

	return result, meta, nil
}
