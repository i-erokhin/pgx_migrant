pgx_migrant
===========

Golang library for migrations and sync interfaces versions of PostgreSQL databases. pgx driver only.

Status
------

It works, but API is unstable. And docs needed.

Usage Example
-------------

```go

package db_sync_helper

import (
	"fmt"

	"github.com/i-erokhin/pgx_migrant"
	"github.com/jackc/pgx"

	"giterica.io/project/repo/config"
	"giterica.io/project/repo/constants"
)

func DBSyncHelper(c *config.Conf, db *pgx.ConnPool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = migrate(c, tx); err != nil {
		return err
	}

	return tx.Commit()
}

func migrate(c *config.Conf, tx *pgx.Tx) error {
	sm := migrant.SyncMap{
		Migration: migrant.MigrationMap{
			Path:   c.SQlArtifacts.MigrationsPath,
			Schema: constants.DBSchema,
		},
		Interfaces: migrant.InterfacesMap{
			Path:   c.SQlArtifacts.InterfacesPath,
			Schema: constants.DBInterfacesSchema,
		},
		TemplateData: c,
	}

	result, err := migrant.SyncDB(tx, sm)
	if err != nil {
		return err
	}
	if len(result.Migrations.Files) != 0 {
		fmt.Printf("DB state: %d, migrations applied:\n", result.Migrations.Number)
		for _, name := range result.Migrations.Files {
			fmt.Printf("  %s\n", name)
		}
	}
	if len(result.Interfaces.Files) != 0 {
		fmt.Printf("Actual DB interfaces schema: %s, SQL files applied:\n", constants.DBInterfacesSchema)
		for _, name := range result.Interfaces.Files {
			fmt.Printf("  %s\n", name)
		}
	}
	return nil
}
```

Output:

```
DB state: 2, migrations applied:
  01-tables.sql
  02-seeds.sql
Actual DB interfaces schema: tabel_interfaces_0, SQL files applied:
  01-active_employee.sql
```