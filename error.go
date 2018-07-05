package migrant

import (
	"encoding/json"
	"fmt"

	"gopkg.in/jackc/pgx.v3"
)

type MigrationError struct {
	File *File
	Err  pgx.PgError
}

func NewError(f *File, err pgx.PgError) *MigrationError {
	return &MigrationError{
		File: f,
		Err:  err,
	}
}

func (m *MigrationError) Error() string {
	b, err := json.MarshalIndent(m.Err, "", "  ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s\nin file: %s\n%s", m.Err.Error(), m.File.Path, b)
}
