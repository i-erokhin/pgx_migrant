package migrant

import (
	"fmt"

	"gopkg.in/jackc/pgx.v3"
)

const TableName = "_migration"

type DB struct {
	Tx         *pgx.Tx
	Schema     string
	LastNumber int
	table      string
}

func NewDB(tx *pgx.Tx, schema string) (db *DB, err error) {
	db = &DB{
		Tx:     tx,
		Schema: schema,
		table:  fmt.Sprintf("%s.%s", schema, TableName),
	}

	var exists bool
	{
		row := tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM information_schema.tables
			WHERE 
				table_schema = $1
          		AND table_name = $2)`,
			db.Schema, TableName)
		if row.Scan(&exists); err != nil {
			return
		}
	}

	if !exists {
		_, err = tx.Exec(fmt.Sprintf(`
			CREATE SCHEMA IF NOT EXISTS %s;
			SET SEARCH_PATH TO %s;
			CREATE TABLE %s (
				number INTEGER      NOT NULL PRIMARY KEY,
				name   VARCHAR(255) NOT NULL UNIQUE,
  				cTime  TIMESTAMP    NOT NULL DEFAULT now()
			);
			LOCK TABLE %s IN ACCESS EXCLUSIVE MODE;`, db.Schema, db.Schema, db.table, db.table))
	} else {
		_, err = tx.Exec(fmt.Sprintf("LOCK TABLE %s IN ACCESS EXCLUSIVE MODE", db.table))
		if err != nil {
			return
		}
		err = tx.QueryRow(fmt.Sprintf("SELECT max(number) FROM %s", db.table)).Scan(&db.LastNumber)
		if err == pgx.ErrNoRows {
			err = nil
		}
	}
	return
}

func (d *DB) Migrate(files []*File, templateData interface{}) (names []string, err error) {
	for _, file := range files {
		if file.Number <= d.LastNumber {
			continue
		}
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
		_, err = d.Tx.Exec(fmt.Sprintf("INSERT INTO %s (number, name) VALUES ($1, $2)", d.table),
			file.Number, file.Basename)
		if err != nil {
			return
		}
		names = append(names, file.Basename)
	}
	return
}
