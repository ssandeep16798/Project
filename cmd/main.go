package main

import (
	"dmarc_backend/internal/engine/gethandler"
	"dmarc_backend/internal/engine/posthandler"
	"dmarc_backend/internal/engine/userhandler"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	//GET API to retrieve data from the DB.
	r.HandleFunc("/getdatabyreportid/", gethandler.HandleGet).Methods("GET")
	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("../internal/infrastructure/swagger/")))
	r.PathPrefix("/swaggerui/").Handler(sh)
	//POST API to insert data into DB.
	r.HandleFunc("/parseXML", posthandler.HandlePost).Methods("POST")
	r.HandleFunc("/login/", userhandler.Login).Methods("POST")
	r.HandleFunc("/registeruser/", userhandler.Register).Methods("POST")
	r.HandleFunc("/verifyemail/", userhandler.VerifyEmail).Methods("POST")
	r.HandleFunc("/resetpassword/", userhandler.ResetPassword).Methods("PUT")
	r.HandleFunc("/getdatabyuserid/", gethandler.GetDataByUserId).Methods("GET")
	r.HandleFunc("/getdatabydomain/", gethandler.HandleGetDataByDomain).Methods("GET")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	server.ListenAndServe()

}
