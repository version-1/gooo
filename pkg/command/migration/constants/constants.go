package constants

const ConfigTableName = "gooo_migration_meta"

type MigrationKind string

const (
	SchemaMigration MigrationKind = "schema"
	DiffMigration   MigrationKind = "diff"
	UpMigration     MigrationKind = "up"
	DownMigration   MigrationKind = "down"
)

type OperationKind string

const (
	AddOperationKind    OperationKind = "add"
	ModifyOperationKind OperationKind = "modify"
	DropOperationKind   OperationKind = "drop"
)
