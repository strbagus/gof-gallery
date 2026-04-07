package helper

import (
	"database/sql"
	"time"
)

func DBToTime(original sql.NullTime) (string, error) {
	var datetime time.Time
	if original.Valid {
		datetime = original.Time
	} else {
		return "", nil
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return "", err
	}
	localized := datetime.In(loc)
	formated := localized.Format(time.RFC3339)

	return formated, nil
}
