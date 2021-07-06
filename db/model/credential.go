package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type Credential struct {
	gorm.Model
	UserID   uint
	Password string
}

func (c *Credential) BeforeCreate(db *gorm.DB) error {
	return c.hashPass(db)
}

func (c *Credential) BeforeUpdate(db *gorm.DB) error {
	if db.Statement.Changed("Password") {
		return c.hashPass(db)
	}
	return nil
}

func (c *Credential) hashPass(db *gorm.DB) error {
	const keyPass = "password"
	var raw string
	switch c := db.Statement.Dest.(type) {
	case map[string]interface{}:
		raw = c[keyPass].(string)
	case *Credential:
		raw = c.Password
	case []*Credential:
		raw = c[db.Statement.CurDestIndex].Password
	default:
		return errors.New("Credential type not supported in preupdate hook!")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err == nil {
		db.Statement.SetColumn(keyPass, hash)
	}
	return err

}

func (c *Credential) isValid(raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(raw))
	return err == nil
}
