package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Date time.Time

func (d *Date) ScanDate(v pgtype.Date) error {
	if !v.Valid {
		return fmt.Errorf("cannot scan NULL into Date")
	}

	if v.InfinityModifier != pgtype.Finite {
		return fmt.Errorf("cannot scan infinite date")
	}

	*d = Date(v.Time)
	return nil
}

func (d Date) DateValue() (pgtype.Date, error) {
	return pgtype.Date{
		Time:             time.Time(d),
		Valid:            true,
		InfinityModifier: pgtype.Finite,
	}, nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

func (d Date) Date() string {
	return time.Time(d).Format("2006-01-02")
}

func (d Date) ViewDate() string {
	return time.Time(d).Format("Jan 2, 2006")
}
