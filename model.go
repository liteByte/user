package main

import (
	"strconv"
	"time"
)

type Model struct {
	ID               uint `gorm:"primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `sql:"index"`
	Email            string     `gorm:"type:varchar(254);unique_index"`
	Password         string     `gorm:"type:char(192)"`
	Compromised      bool
	ProtectionScheme string `gorm:"type:char(32)"`
	Name             string
	Age              uint
	Number           int
	Date             time.Time
}

func ParseAgeFromString(s string) (uint, error) {
	age, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(age), nil
}

func ParseNumberFromString(s string) (int, error) {
	number, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(number), nil
}

func ParseDateFromString(s string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, s)
	if err != nil {
		println(err.Error())
		return time.Time{}, err
	}
	return date, nil
}
