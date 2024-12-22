package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"quiz-service/init/logger"
	"quiz-service/internal/entities"
)

type Testing struct {
	db     *sqlx.DB
	logger logger.Logging
}

func NewTesting(db *sqlx.DB, logger logger.Logging) *Testing {
	return &Testing{db: db, logger: logger}
}

func (t *Testing) TestGet(ctx context.Context, userId, variantId int) (*entities.Testing, error) {
	t.logger.InfoF("TestGet received | %d | %d", userId, variantId)

	var testEntity = new(entities.Testing)
	query := `
		SELECT id, user_id, variant_id, correct_answers, start_at, finish_at 
		FROM testing 
		WHERE user_id = $1 AND variant_id = $2
	`
	if err := t.db.GetContext(ctx, testEntity, query, userId, variantId); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	t.logger.InfoF("TestGet success | %d | %d", userId, variantId)

	return testEntity, nil
}
