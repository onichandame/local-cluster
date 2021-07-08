package model

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EnumValue string

const (
	PENDING     EnumValue = "PENDING"
	CREATING    EnumValue = "CREATING"
	RUNNING     EnumValue = "RUNNING"
	TERMINATING EnumValue = "TERMINATING"
	FINISHED    EnumValue = "FINISHED"
	FAILED      EnumValue = "FAILED"
	NEVER       EnumValue = "NEVER"
	ALWAYS      EnumValue = "ALWAYS"
)

type Enum struct {
	gorm.Model
	Value EnumValue `gorm:"unique"`
}

func getEnums(db *gorm.DB) (map[EnumValue]*Enum, error) {
	recChan := make(chan *Enum)
	errChan := make(chan error)
	getRec := func(status EnumValue) {
		rec := Enum{}
		err := db.FirstOrCreate(&rec, Enum{Value: status}).Error
		if err != nil {
			errChan <- err
		} else {
			recChan <- &rec
		}
	}
	statuses := []EnumValue{PENDING, FINISHED, FAILED, RUNNING, TERMINATING, CREATING, ALWAYS, NEVER}
	for _, s := range statuses {
		go getRec(s)
	}
	res := make(map[EnumValue]*Enum)
	for range statuses {
		select {
		case err := <-errChan:
			return nil, err
		case rec := <-recChan:
			res[rec.Value] = rec
		}
	}
	return res, nil
}

func selectEnums(db *gorm.DB, enumValues []EnumValue) map[EnumValue]*Enum {
	evMaps := map[EnumValue]interface{}{}
	for _, e := range enumValues {
		evMaps[e] = nil
	}
	rec, err := getEnums(db)
	if err != nil {
		logrus.Fatalf("failed to get enum")
	}
	for key := range rec {
		_, ok := evMaps[key]
		if !ok {
			delete(rec, key)
		}
	}
	return rec
}

func GetJobStatuses(db *gorm.DB) map[EnumValue]*Enum {
	allowedEnumValues := []EnumValue{PENDING, FINISHED, FAILED}
	return selectEnums(db, allowedEnumValues)
}

func GetInstanceStatuses(db *gorm.DB) map[EnumValue]*Enum {
	allowedEnumValues := []EnumValue{PENDING, CREATING, RUNNING, TERMINATING, FAILED, FINISHED}
	return selectEnums(db, allowedEnumValues)
}

func GetRestartPolicies(db *gorm.DB) map[EnumValue]*Enum {
	allowedEnumValues := []EnumValue{ALWAYS, NEVER}
	return selectEnums(db, allowedEnumValues)
}
