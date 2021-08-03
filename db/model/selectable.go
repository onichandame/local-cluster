package model

type Selectable struct {
	Name string `gorm:"unique,not null"`
}

type selectable struct {
	Name string
}
