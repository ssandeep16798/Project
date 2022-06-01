package models

type Response struct {
	Status   bool
	Msg      string
	Data     []Data
	ReportID string
}
