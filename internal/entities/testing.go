package entities

import "time"

type Testing struct {
	ID             int        `json:"id"`
	UserId         int        `json:"user_id" db:"user_id"`
	VariantId      int        `json:"variant_id" db:"variant_id"`
	CorrectAnswers int        `json:"correct_answers" db:"correct_answers"`
	StartAt        time.Time  `json:"start_at" db:"start_at"`
	FinishAt       *time.Time `json:"finish_at" db:"finish_at"`
}
