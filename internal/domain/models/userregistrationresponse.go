package models

type UserRegisterResponse struct {
	Status bool
	Msg    string
	Data   RegisterResponse
}
