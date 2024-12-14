package schema

import (
	"fmt"

	"github.com/version-1/gooo/pkg/command/migration/adapter/yaml"
)

type MigrationConfig struct {
	TableNameMapper map[string]string
	Indexes         map[string][]yaml.Index
}

func NewMigration(collection SchemaCollection, config MigrationConfig) *Migration {
	m := Migration{
		collection: collection,
		config:     config,
	}

	if m.config.Indexes == nil {
		m.config.Indexes = map[string][]yaml.Index{}
	}

	return &m
}

type Migration struct {
	collection SchemaCollection
	config     MigrationConfig
}

func (m Migration) OriginSchema() (yaml.OriginSchema, error) {
	schema := yaml.OriginSchema{}
	for _, s := range m.collection.Schemas {
		columns := []yaml.Column{}
		for _, f := range s.Fields {
			if f.IsAssociation() {
				continue
			}

			columns = append(columns, yaml.Column{
				Name:       f.ColumnName(),
				Type:       f.TableType(),
				Default:    &f.Tag.DefaultValue,
				AllowNull:  &f.Tag.AllowNull,
				PrimaryKey: &f.Tag.PrimaryKey,
			})
		}

		indexes := m.config.Indexes[s.Name]
		for _, f := range s.Fields {
			if !f.IsAssociation() && (f.Tag.Index || f.Tag.Unique) {
				indexes = append(indexes, yaml.Index{
					Name:    fmt.Sprintf("index_%s_%s", s.TableName, f.ColumnName()),
					Columns: []string{f.ColumnName()},
					Unique:  &f.Tag.Unique,
				})
			}
		}

		tableName, ok := m.config.TableNameMapper[s.Name]
		if !ok {
			tableName = s.TableName
		}
		schema.Tables = append(schema.Tables, yaml.Table{
			Name:    tableName,
			Columns: columns,
			Indexes: indexes,
		})
	}

	return schema, nil
}
