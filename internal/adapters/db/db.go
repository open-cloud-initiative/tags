package db

import (
	"github.com/open-cloud-initiative/tags/internal/ports"

	"github.com/katallaxie/pkg/dbx"
	"gorm.io/gorm"
)

var _ ports.ReadTx = (*readTxImpl)(nil)

type readTxImpl struct {
	conn *gorm.DB
}

// NewReadTx ...
func NewReadTx() dbx.ReadTxFactory[ports.ReadTx] {
	return func(db *gorm.DB) (ports.ReadTx, error) {
		return &readTxImpl{conn: db}, nil
	}
}

type writeTxImpl struct {
	conn *gorm.DB
	readTxImpl
}

// NewWriteTx ...
func NewWriteTx() dbx.ReadWriteTxFactory[ports.ReadWriteTx] {
	return func(db *gorm.DB) (ports.ReadWriteTx, error) {
		return &writeTxImpl{conn: db}, nil
	}
}
