package model

type Selectable struct {
	Name string `gorm:"unique"`
}
