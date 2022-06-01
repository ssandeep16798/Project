package models

type LoginResponse struct {
	Status bool
	Msg    string
}

type LoginResponseWithData struct {
	Status bool
	Msg    string
	Data   UserRegistrationRequest
}
