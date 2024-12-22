package entities

import "time"

type Register struct {
	UUID     string `json:"uuid"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Login struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID           int        `json:"id"`
	UUID         string     `json:"uuid"`
	Login        string     `json:"login"`
	Authorized   bool       `json:"authorized"`
	AuthorizedAt time.Time  `json:"authorized_at" db:"authorized_at"`
	QuitAt       *time.Time `json:"quit_at,omitempty" db:"quit_at"`
}
