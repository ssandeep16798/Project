package models

import "time"

type UserRegistrationRequest struct {
	Id               int
	Firstname        string
	Lastname         string
	Username         string
	Email            string
	Password         string
	Gender           string
	Mobile           string
	Created_datetime time.Time
}
