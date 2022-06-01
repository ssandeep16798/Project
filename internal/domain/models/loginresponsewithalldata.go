package models

type LoginResponseWithAllData struct {
	Status bool
	Msg    string
	Data   []UserRegistrationRequest
}
