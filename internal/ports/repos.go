package ports

import (
	"context"

	"github.com/katallaxie/pkg/dbx"
)

// Datastore provides methods for transactional operations.
type Datastore interface {
	// ReadTx starts a read only transaction.
	ReadTx(context.Context, func(context.Context, ReadTx) error) error
	// ReadWriteTx starts a read write transaction.
	ReadWriteTx(context.Context, func(context.Context, ReadWriteTx) error) error

	dbx.Migrator
}

// ReadTx provides methods for transactional read operations.
type ReadTx interface{}

// ReadWriteTx provides methods for transactional read and write operations.
type ReadWriteTx interface {
	ReadTx
}
