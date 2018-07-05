package migrant

import (
	"fmt"

	"gopkg.in/jackc/pgx.v3"
)

type DBInterfaces struct {
	Tx     *pgx.Tx
	Schema string
	Found  bool
}

func NewDbInterfaces(tx *pgx.Tx, schema string) (result *DBInterfaces, err error) {
	result = &DBInterfaces{
		Tx:     tx,
		Schema: schema,
	}
	row := tx.QueryRow(`
		SELECT exists(
		    SELECT schema_name
		    FROM information_schema.schemata
		    WHERE schema_name = $1)`, schema)
	err = row.Scan(&result.Found)
	return
}

func (d *DBInterfaces) Sync(files []*File, templateData interface{}) (names []string, err error) {
	if d.Found {
		panic("wat?")
	}

	_, err = d.Tx.Exec(fmt.Sprintf(`
		CREATE SCHEMA %s;
		SET SEARCH_PATH TO %s;`, d.Schema, d.Schema))
	if err != nil {
		return
	}

	for _, file := range files {
		var sql string
		sql, err = file.GetSQL(templateData)
		if err != nil {
			return
		}
		_, err = d.Tx.Exec(sql)
		if err != nil {
			if pgErr, ok := err.(pgx.PgError); ok {
				err = NewError(file, pgErr)
			}
			return
		}
		names = append(names, file.Basename)
	}
	return
}