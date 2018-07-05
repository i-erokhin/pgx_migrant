package migrations

import (
	"gopkg.in/jackc/pgx.v3"
)

type MigrationMap struct {
	Schema string
	Path   string
}

type InterfacesMap MigrationMap

type SyncMap struct {
	Migration    MigrationMap
	Interfaces   InterfacesMap
	TemplateData interface{}
}

type MigrationResult struct {
	Files  []string
	Number int
}

type InterfacesResult struct {
	Files []string
}

type InitDBResult struct {
	Migrations MigrationResult
	Interfaces InterfacesResult
}

func SyncDB(tx *pgx.Tx, sm SyncMap) (result InitDBResult, err error) {
	{
		var mr MigrationResult
		mr, err = Migrate(tx, sm.Migration, sm.TemplateData)
		if err != nil {
			return
		}
		result.Migrations = mr
	}

	{
		var ir InterfacesResult
		ir, err = CreateInterfaces(tx, sm.Interfaces, sm.TemplateData)
		if err != nil {
			return
		}
		result.Interfaces = ir
	}
	return
}


func Migrate(tx *pgx.Tx, mm MigrationMap, templateData interface{}) (result MigrationResult, err error) {
	fs, err := NewFilesystem(mm.Path)
	if err != nil {
		return
	}
	dbMigrator, err := NewDB(tx, mm.Schema)
	if err != nil {
		return
	}
	result.Number = fs.MaxNumber
	result.Files, err = dbMigrator.Migrate(fs.Files, templateData)
	return
}

func CreateInterfaces(tx *pgx.Tx, im InterfacesMap, templateData interface{}) (result InterfacesResult, err error) {
	interfacer, err := NewDbInterfaces(tx, im.Schema)
	if err != nil || interfacer.Found {
		return
	}
	fs, err := NewFilesystem(im.Path)
	if err != nil {
		return
	}
	result.Files, err = interfacer.Sync(fs.Files, templateData)
	return
}
