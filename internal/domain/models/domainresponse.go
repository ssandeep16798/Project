package models

type DomainResponse struct {
	Status     bool
	Msg        string
	DomainName string
	Data       []DomainData
}
