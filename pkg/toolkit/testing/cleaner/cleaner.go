package cleaner

import (
	"context"

	"github.com/version-1/gooo/pkg/db"
	"github.com/version-1/gooo/pkg/testing/cleaner/adapter"
)

var _ CleanAdapter = (*adapter.Pq)(nil)

type CleanAdapter interface {
	ListTables(ctx context.Context) ([]string, error)
	Truncate(ctx context.Context, table string) error
	ResetIndexes(ctx context.Context, table string) error
}

type Cleaner struct {
	adapter CleanAdapter
}

func New(conn db.Tx) *Cleaner {
	return &Cleaner{adapter: adapter.New(conn)}
}

func (c Cleaner) Clean(ctx context.Context) {
	tables, err := c.adapter.ListTables(ctx)
	if err != nil {
		panic(err)
	}

	for _, table := range tables {
		if err := c.adapter.Truncate(ctx, table); err != nil {
			panic(err)
		}

		// INFO: have to reset index after truncate mainly for pkey and unique index
		if err := c.adapter.ResetIndexes(ctx, table); err != nil {
			panic(err)
		}
	}
}
