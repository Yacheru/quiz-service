package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jmoiron/sqlx"
	"quiz-service/init/logger"
	"quiz-service/internal/entities"
	"quiz-service/pkg/constants"
	"time"
)

type Variant struct {
	db     *sqlx.DB
	logger logger.Logging
}

func NewVariant(db *sqlx.DB, logger logger.Logging) *Variant {
	return &Variant{db: db, logger: logger}
}

func (v *Variant) VariantAdd(ctx context.Context, name string) error {
	v.logger.InfoF("VariantAdd received | %s", name)

	query := `
		INSERT INTO variants (name) VALUES ($1);
	`
	if _, err := v.db.ExecContext(ctx, query, name); err != nil {
		return err
	}

	v.logger.InfoF("VariantAdd success | %s", name)

	return nil
}

func (v *Variant) VariantRemove(ctx context.Context, name string) (int64, error) {
	v.logger.InfoF("VariantRemove received | %s", name)

	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	variantQuery := `
		DELETE FROM variants WHERE name = $1;
	`
	result, err := tx.ExecContext(ctx, variantQuery, name)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	answersQuery := `
		DELETE FROM answers
		WHERE id NOT IN (
			SELECT DISTINCT answers_id
			FROM questions_and_answers
		);
	`
	if _, err := tx.ExecContext(ctx, answersQuery); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	v.logger.InfoF("VariantRemove success | %s", name)

	return result.RowsAffected()
}

func (v *Variant) VariantList(ctx context.Context) ([]*entities.Variant, error) {
	v.logger.Info("VariantRemove received")

	query := `
		SELECT
			v.id, v.name,
			q.id AS question_id, q.question, q.answer,
			(
				SELECT json_agg(json_build_object('answer', a.answer))
				FROM questions_and_answers qaa
				LEFT JOIN public.answers a ON a.id = qaa.answers_id
				WHERE qaa.questions_id = q.id
			) AS answers
		FROM variants v
			LEFT JOIN questions q ON v.id = q.variant_id
		ORDER BY v.id, q.id;
	`

	rows, err := v.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants = make([]*entities.Variant, 0)
	var currentVariant *entities.Variant

	for rows.Next() {
		var (
			variantId      sql.Null[int]
			questionId     sql.Null[int]
			variantName    sql.Null[string]
			questionName   sql.Null[string]
			questionAnswer sql.Null[string]
			answersByte    []byte
		)

		if err := rows.Scan(&variantId, &variantName, &questionId, &questionName, &questionAnswer, &answersByte); err != nil {
			return nil, err
		}

		if currentVariant == nil || currentVariant.Id != variantId.V {
			if currentVariant != nil {
				variants = append(variants, currentVariant)
			}

			currentVariant = &entities.Variant{
				Id:        variantId.V,
				Name:      variantName.V,
				Questions: make([]*entities.Question, 0),
			}
		}

		if len(answersByte) == 0 {
			continue
		}

		var answersArr []*entities.Answer
		if err := json.Unmarshal(answersByte, &answersArr); err != nil {
			return nil, err
		}

		currentVariant.Questions = append(currentVariant.Questions, &entities.Question{
			Id:       questionId.V,
			Question: questionName.V,
			Answer:   questionAnswer.V,
			Answers:  answersArr,
		})
	}

	if currentVariant != nil {
		variants = append(variants, currentVariant)
	}

	v.logger.Info("VariantRemove received")

	return variants, nil
}

func (v *Variant) VariantGet(ctx context.Context, name string) (*entities.Variant, error) {
	v.logger.InfoF("VariantGet received | %s", name)

	query := `
		SELECT
			v.id, v.name,
			q.id AS question_id, q.question, q.answer,
			(
				SELECT json_agg(json_build_object('answer', a.answer))
				FROM questions_and_answers qaa
				LEFT JOIN answers a ON a.id = qaa.answers_id
				WHERE qaa.questions_id = q.id
			) AS answers
		FROM variants v
			LEFT JOIN questions q ON v.id = q.variant_id
		WHERE v.name = $1;
	`

	rows, err := v.db.QueryxContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questionsMap := make(map[int]*entities.Question)

	var variantEntity = new(entities.Variant)
	for rows.Next() {
		var (
			variantId   sql.Null[int]
			variantName sql.Null[string]
			questionId  sql.Null[int]
			question    sql.Null[string]
			answer      sql.Null[string]
			answers     []byte
		)

		if err := rows.Scan(&variantId, &variantName, &questionId, &question, &answer, &answers); err != nil {
			return nil, err
		}

		if variantEntity.Id == 0 && variantId.Valid {
			variantEntity.Id = variantId.V
		}
		if variantEntity.Name == "" && variantName.Valid {
			variantEntity.Name = variantName.V
		}

		if questionId.Valid {
			qId := questionId.V
			if _, exists := questionsMap[qId]; !exists {
				questionsMap[qId] = &entities.Question{
					Id:       qId,
					Question: question.V,
					Answer:   answer.V,
				}
				variantEntity.Questions = append(variantEntity.Questions, questionsMap[qId])
			}

			if answers != nil {
				var parsedAnswers []map[string]string
				if err := json.Unmarshal(answers, &parsedAnswers); err != nil {
					return nil, err
				}

				for _, a := range parsedAnswers {
					if ans, ok := a["answer"]; ok {
						questionsMap[qId].Answers = append(questionsMap[qId].Answers, &entities.Answer{Answer: ans})
					}
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(variantEntity.Questions) == 0 {
		variantEntity.Questions = make([]*entities.Question, 0)
	}

	if variantEntity.Id == 0 {
		return nil, constants.ErrorVariantNotFound
	}

	v.logger.InfoF("VariantGet success | %s", name)

	return variantEntity, nil
}

func (v *Variant) VariantStart(ctx context.Context, variantId, userId int) error {
	v.logger.InfoF("VariantStart received | %d | %d", variantId, userId)

	var finish *time.Time
	selectQuery := `
		SELECT finish_at FROM testing
		WHERE user_id = $1 AND variant_id = $2;
	`
	if err := v.db.GetContext(ctx, &finish, selectQuery, userId, variantId); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if finish != nil {
		return constants.ErrorVariantCompleted
	}

	query := `
		INSERT INTO testing (user_id, variant_id) VALUES ($1, $2);
	`
	if _, err := v.db.ExecContext(ctx, query, userId, variantId); err != nil {
		return err
	}

	v.logger.InfoF("VariantStart success | %d | %d", variantId, userId)

	return nil
}

func (v *Variant) VariantResults(ctx context.Context, variantId, userId int) (*entities.Testing, error) {
	v.logger.InfoF("VariantResults received | %d | %d", variantId, userId)

	tx, err := v.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	testingEntity := new(entities.Testing)
	finishTestingQuery := `
		UPDATE testing SET finish_at = $1 
		WHERE user_id = $2 AND variant_id = $3
		RETURNING id, user_id, variant_id, correct_answers, start_at, finish_at;
	`
	if err := v.db.GetContext(ctx, testingEntity, finishTestingQuery, now, userId, variantId); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	v.logger.InfoF("VariantResults success | %d | %d", variantId, userId)

	return testingEntity, nil
}
