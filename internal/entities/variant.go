package entities

type Variant struct {
	Id        int         `json:"id"`
	Name      string      `json:"name" binding:"required"`
	Questions []*Question `json:"questions"`
}

type Results struct {
	Percent int `json:"percent"`
}
