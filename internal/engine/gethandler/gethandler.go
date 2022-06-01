package gethandler

import (
	"dmarc_backend/internal/domain/models"
	"dmarc_backend/internal/engine/DB"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func HandleGet(w http.ResponseWriter, req *http.Request) {

	fmt.Println("Inside Handle Get")
	queryParams := req.URL.Query()
	fmt.Println(queryParams)

	w.Header().Set("Content-Type", "application/json")
	reportId := ""
	domainName := ""
	if queryParams["report_id"] != nil && len(queryParams["report_id"]) > 0 {

		reportId = queryParams["report_id"][0]
	}
	if queryParams["domain_name"] != nil && len(queryParams["domain_name"]) > 0 {
		domainName = queryParams["domain_name"][0]
	}

	if reportId == "" && domainName == "" {
		var response models.Response

		response.Status = false
		response.Msg = "Fields Required"
		response.Data = make([]models.Data, 0)

		json.NewEncoder(w).Encode(response)
		return

	}

	var data []models.Data = DB.GetDmarcFeedbackData(reportId, domainName)
	response := HandleResponse(data)

	//feedbackArr := make([]models.Feedback, len(data))

	/*for i := 0; i < len(data); i++ {
		json.Unmarshal([]byte(data[i]), &feedbackArr[i])
	}*/

	json.NewEncoder(w).Encode(response)
	//w.Write(resp)

}
func HandleResponse(data []models.Data) models.Response {
	var response models.Response
	if len(data) >= 1 {
		response.Data = data
		response.Msg = "Data Found"
		response.Status = true
	} else {
		response.Data = data
		response.Msg = "Data Not Found"
		response.Status = false
	}
	return response
}

func GetDataByUserId(w http.ResponseWriter, req *http.Request) {

	fmt.Println("Inside Get Data By User ID")
	queryParams := req.URL.Query()
	fmt.Println(queryParams)

	w.Header().Set("Content-Type", "application/json")

	userid := ""

	if queryParams["user_id"] != nil && len(queryParams["user_id"]) > 0 {

		userid = queryParams["user_id"][0]
	}
	if userid == "" {
		var response models.Response

		response.Status = false
		response.Msg = "userID Required"
		response.Data = make([]models.Data, 0)

		json.NewEncoder(w).Encode(response)
		return

	}

	useridint, _ := strconv.Atoi(userid)
	var data models.UserId = DB.GetDmarcFeedbackDataByUserId(useridint)

	json.NewEncoder(w).Encode(data)

}

func HandleGetDataByDomain(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inside Handle Get Data by domain")
	queryParams := req.URL.Query()
	start_date := ""
	end_date := ""
	fmt.Println(queryParams)
	var response models.DomainResponse

	w.Header().Set("Content-Type", "application/json")

	domainName := ""

	if queryParams["domain_name"] != nil && len(queryParams["domain_name"]) > 0 {
		domainName = queryParams["domain_name"][0]
	}

	if queryParams["start_date"] != nil && len(queryParams["start_date"]) > 0 {
		start_date = queryParams["start_date"][0]
	}

	if queryParams["end_date"] != nil && len(queryParams["end_date"]) > 0 {
		end_date = queryParams["end_date"][0]
	}

	fmt.Println(domainName, start_date, end_date)

	if domainName == "" || start_date == "" || end_date == "" {

		response.Status = false
		response.Msg = "Fields Required"

		json.NewEncoder(w).Encode(response)
		return

	}
	data := DB.GetDmarcFeedBackDataByDomain(domainName, start_date, end_date)

	data.DomainName = domainName

	json.NewEncoder(w).Encode(data)
}
