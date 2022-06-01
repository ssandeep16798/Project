package models

import "time"

type Users struct {
	Id               int64
	Report_id        string
	Domain_name      string
	Username         string
	Created_datetime time.Time
	User_id          int64
}
