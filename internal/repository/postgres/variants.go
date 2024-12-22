package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Variant struct {
	db *sqlx.DB
}

func NewVariant(db *sqlx.DB) *Variant {
	return &Variant{db: db}
}

func (v *Variant) VariantAdd(ctx context.Context, name string) error {
	query := `
		INSERT INTO variants (name) VALUES ($1);
	`
	if _, err := v.db.ExecContext(ctx, query, name); err != nil {
		return err
	}
	return nil
}

func (v *Variant) VariantRemove(ctx context.Context, name string) (int64, error) {
	query := `
		DELETE FROM variants WHERE name = $1;
	`
	result, err := v.db.ExecContext(ctx, query, name)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
