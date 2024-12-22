package entities

type Answer struct {
	Answer string `json:"answer" binding:"required"`
}
