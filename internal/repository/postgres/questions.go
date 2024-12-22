package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"quiz-service/init/logger"
	"quiz-service/internal/entities"
	"sync"
)

type Questions struct {
	db     *sqlx.DB
	logger logger.Logging

	mu sync.Mutex
}

func NewQuestions(db *sqlx.DB, logger logger.Logging) *Questions {
	return &Questions{db: db, mu: sync.Mutex{}, logger: logger}
}

func (q *Questions) QuestionCount(ctx context.Context, variantId int) (int, error) {
	q.logger.InfoF("QuestionCount received | %d", variantId)

	var count int
	query := `
		SELECT COUNT(*) FROM questions WHERE variant_id = $1
	`
	if err := q.db.GetContext(ctx, &count, query, variantId); err != nil {
		return 0, err
	}

	q.logger.InfoF("QuestionCount success | %d", variantId)

	return count, nil
}

func (q *Questions) QuestionAdd(ctx context.Context, variantId int, question *entities.Question) error {
	q.logger.InfoF("QuestionAdd received %d | %+v", variantId, question)

	tx, err := q.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	var questionId int
	questionQuery := `
		INSERT INTO questions (variant_id, question, answer) VALUES ($1, $2, $3) RETURNING id;
	`
	if err := tx.GetContext(ctx, &questionId, questionQuery, variantId, question.Question, question.Answer); err != nil {
		tx.Rollback()
		return err
	}

	var answerIds = make([]int, 0, len(question.Answers))
	for _, answer := range question.Answers {
		var answerId int
		answerQuery := `
			INSERT INTO answers (answer) VALUES ($1) RETURNING id;
		`
		if err := tx.GetContext(ctx, &answerId, answerQuery, answer.Answer); err != nil {
			tx.Rollback()
			return err
		}
		answerIds = append(answerIds, answerId)
	}

	for _, answerId := range answerIds {
		questionsAndAnswersQuery := `
			INSERT INTO questions_and_answers (questions_id, answers_id) VALUES ($1, $2)
		`
		if _, err := tx.ExecContext(ctx, questionsAndAnswersQuery, questionId, answerId); err != nil {
			tx.Rollback()
			return err
		}
	}

	q.logger.InfoF("QuestionAdd success | %d | %+v", variantId, question)

	return tx.Commit()
}

func (q *Questions) QuestionRemove(ctx context.Context, variantId int, question string) (int64, error) {
	q.logger.InfoF("QuestionRemove received %d | %s", variantId, question)

	tx, err := q.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return 0, err
	}

	queryDeleteLinks := `
		DELETE FROM questions_and_answers
		WHERE questions_id = (SELECT id FROM questions WHERE variant_id = $1 AND question = $2)
	`
	_, err = tx.ExecContext(ctx, queryDeleteLinks, variantId, question)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	queryDeleteAnswers := `
		DELETE FROM answers
		WHERE id NOT IN (SELECT answers_id FROM questions_and_answers)
	`
	_, err = tx.ExecContext(ctx, queryDeleteAnswers)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	queryDeleteQuestion := `
		DELETE FROM questions WHERE variant_id = $1 AND question = $2
	`
	res, err := tx.ExecContext(ctx, queryDeleteQuestion, variantId, question)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	q.logger.InfoF("QuestionRemove success %d | %s", variantId, question)

	return res.RowsAffected()
}

func (q *Questions) QuestionGet(ctx context.Context, variantId, questionId int) (*entities.Question, error) {
	q.logger.InfoF("QuestionGet received %d | %d", variantId, questionId)

	var question = new(entities.Question)
	var answers = new([]byte)

	query := `
		SELECT
			q.id,
			q.question,
			q.answer,
			json_agg(json_build_object('answer', ans.answer)) AS answers
		FROM questions q
			JOIN variants v ON v.id = q.variant_id
			JOIN questions_and_answers qa ON qa.questions_id = q.id
			JOIN answers ans ON ans.id = qa.answers_id
		WHERE q.id = $1 AND q.variant_id = $2
		GROUP BY q.question, v.name, q.answer, q.id
	`
	if err := q.db.QueryRowxContext(ctx, query, questionId, variantId).Scan(&question.Id, &question.Question, &question.Answer, answers); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(*answers, &question.Answers); err != nil {
		return nil, err
	}

	q.logger.InfoF("QuestionGet success %d | %d", variantId, questionId)

	return question, nil
}

func (q *Questions) QuestionAccept(ctx context.Context, answer string, userId, testId, variantId int) error {
	q.logger.InfoF("QuestionGet received %s | %d | %d | %d", variantId, userId, testId, answer)

	var questionEntity = new(entities.Question)
	query := `
		SELECT id, answer, question
		FROM questions 
		WHERE variant_id = $1
	`
	if err := q.db.GetContext(ctx, questionEntity, query, variantId); err != nil {
		return err
	}

	insertQuery := `
			INSERT INTO user_answers (test_id, question_id, answer) VALUES ($1, $2, $3)
	`
	if questionEntity.Answer == answer {
		q.mu.Lock()
		if _, err := q.db.ExecContext(ctx, insertQuery, testId, questionEntity.Id, answer); err != nil {
			return err
		}

		updateQuery := `
			UPDATE testing SET correct_answers = correct_answers + 1
			WHERE user_id = $1 AND variant_id = $2
		`
		if _, err := q.db.ExecContext(ctx, updateQuery, userId, variantId); err != nil {
			return err
		}

		q.mu.Unlock()
	}

	if _, err := q.db.ExecContext(ctx, insertQuery, testId, questionEntity.Id, answer); err != nil {
		return err
	}

	q.logger.InfoF("QuestionGet success %s | %d | %d | %d", variantId, userId, testId, answer)

	return nil
}
