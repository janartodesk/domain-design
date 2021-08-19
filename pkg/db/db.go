package db

import (
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

// WithTransaction runs a callback within a database transaction.
func WithTransaction(db *bun.DB, callback func(bun.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := callback(tx); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			errors.Wrap(err, txErr.Error())
		}

		return err
	}

	return tx.Commit()
}
