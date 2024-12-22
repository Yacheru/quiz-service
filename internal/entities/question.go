package entities

import "encoding/json"

type Question struct {
	Id       int       `json:"id" db:"id"`
	Question string    `json:"question" binding:"required,max=50" db:"question"`
	Answer   string    `json:"answer" binding:"required,max=50" db:"answer"`
	Answers  []*Answer `json:"answers" binding:"required,len=3"`
}

type QuestionRemove struct {
	VariantName string `json:"variant_name"`
	Question    string `json:"question" binding:"required,max=50"`
}

type QuestionMarshal struct {
	VariantName string          `json:"variant_name" db:"variant_name"`
	Question    string          `json:"question" binding:"required,max=50"`
	Answer      string          `json:"answer" binding:"required,max=50"`
	Answers     json.RawMessage `json:"answers" binding:"required,len=3"`
}
