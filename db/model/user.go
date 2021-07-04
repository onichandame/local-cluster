package model

import "gorm.io/gorm"

type Role int

const (
	ADMIN      Role = iota
	MAINTAINER Role = iota
	GUEST      Role = iota
)

type User struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Credentials []Credential
	Role        Role
}
