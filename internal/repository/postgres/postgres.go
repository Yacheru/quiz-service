package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"quiz-service/init/config"
	"quiz-service/init/logger"
)

func InitPostgresConnection(ctx context.Context, cfg *config.Config, logger logger.Logging) (*sqlx.DB, error) {
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database)

	db, err := sqlx.ConnectContext(ctx, "pgx", uri)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	m, err := migrate.New("file://./migrations", uri)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Debug("migrations already up to date")
		} else {
			logger.Error(err.Error())
			return nil, err
		}
	}

	return db, nil
}
