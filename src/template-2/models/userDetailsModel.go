package models

import "time"

type UserDetails struct {
	ID          int32
	FirstName   string
	LastName    string
	Email       string
	Password    string
	BirthDate   time.Time
	CreatedDate time.Time
}
