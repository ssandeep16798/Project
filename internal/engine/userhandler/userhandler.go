package userhandler

import (
	"dmarc_backend/internal/domain/models"
	"dmarc_backend/internal/engine/DB"
	"encoding/json"
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, req *http.Request) {

	var LoginResponse models.LoginResponse

	fmt.Println("Inside User Login")
	decoder := json.NewDecoder(req.Body)
	var loginRequest models.LoginRequest
	decoder.Decode(&loginRequest)

	w.Header().Set("Content-Type", "application/json")

	var response models.UserRegistrationRequest

	if loginRequest.Username == "" {
		LoginResponse.Status = false
		LoginResponse.Msg = "Fields Required"
		json.NewEncoder(w).Encode(LoginResponse)

	} else {
		count := DB.UserDB(loginRequest.Username, loginRequest.Password)
		isAdmin := DB.CheckIfAdmin(loginRequest.Username, loginRequest.Password)
		fmt.Println(isAdmin)

		if count > 0 {
			var LoginResponseWithData models.LoginResponseWithData
			LoginResponseWithData.Status = true
			LoginResponseWithData.Msg = ""

			response = DB.FetchData(loginRequest.Username, loginRequest.Password)
			LoginResponseWithData.Data = response
			json.NewEncoder(w).Encode(LoginResponseWithData)
		} else {
			LoginResponse.Status = false
			LoginResponse.Msg = "Invalid Login Details"
			json.NewEncoder(w).Encode(LoginResponse)
		}

	}

}

func Register(w http.ResponseWriter, req *http.Request) {

	var LoginResponse models.LoginResponse

	fmt.Println("Inside Register New User")
	decoder := json.NewDecoder(req.Body)
	var userRegistration models.UserRegistrationRequest
	decoder.Decode(&userRegistration)

	w.Header().Set("Content-Type", "application/json")

	count := DB.CheckIfEmailExists(userRegistration.Email, userRegistration.Username)

	if count > 0 {

		LoginResponse.Status = false
		LoginResponse.Msg = "Email or username already Exists"
		json.NewEncoder(w).Encode(LoginResponse)

	} else {
		DB.ResgisterUser(userRegistration)
		fmt.Println("In Else")
		var LoginResponseWithData models.LoginResponseWithData
		LoginResponseWithData.Status = true
		LoginResponseWithData.Msg = ""
		response := DB.FetchData(userRegistration.Username, userRegistration.Password)
		LoginResponseWithData.Data = response
		fmt.Println(response.Email)
		json.NewEncoder(w).Encode(LoginResponseWithData)
	}

}

func ResetPassword(w http.ResponseWriter, req *http.Request) {

	var LoginResponse models.LoginResponse

	decoder := json.NewDecoder(req.Body)
	var userRegistration models.UserRegistrationRequest
	decoder.Decode(&userRegistration)

	w.Header().Set("Content-Type", "application/json")

	DB.UpdatePassword(userRegistration.Email, userRegistration.Password)

	LoginResponse.Status = true
	LoginResponse.Msg = "Password updated successfully."

	fmt.Println("Inside Reset Password")

}

func VerifyEmail(w http.ResponseWriter, request *http.Request) {

	var LoginResponse models.LoginResponse
	decoder := json.NewDecoder(request.Body)
	var userRegistration models.UserRegistrationRequest
	decoder.Decode(&userRegistration)

	w.Header().Set("Content-Type", "application/json")

	count := DB.EmailVerification(userRegistration.Email)

	if count < 1 || userRegistration.Email == "" {

		LoginResponse.Status = false
		LoginResponse.Msg = "Invalid Email"

	} else {
		DB.ResgisterUser(userRegistration)
		LoginResponse.Status = true
		LoginResponse.Msg = "Account exists"

	}
	json.NewEncoder(w).Encode(LoginResponse)

}
