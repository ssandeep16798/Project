package posthandler

import (
	"bytes"
	"dmarc_backend/internal/domain/models"
	"dmarc_backend/internal/engine/structconverter"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"dmarc_backend/internal/engine/DB"
)

func HandlePost(w http.ResponseWriter, req *http.Request) {

	userId, err := strconv.Atoi(req.Form.Get("user_id"))
	if err != nil {
		userId = 1
	}
	file, _, err := req.FormFile("xmlFile")
	fmt.Println(userId)

	var response models.Response
	w.Header().Set("Content-Type", "application/json")

	response.Status = false
	response.Data = make([]models.Data, 0)
	response.Msg = "Fields Required"
	if err != nil {

		fmt.Println("err in handling post", err.Error())
		json.NewEncoder(w).Encode(response)
		return

	}

	if file == nil {
		fmt.Println("File is Required")
		json.NewEncoder(w).Encode(response)
		return
	}
	xml := bytes.NewBuffer(nil)
	io.Copy(xml, file)
	//xml, _ := ioutil.ReadAll(req.Body)
	feedback := structconverter.ConvertToStruct(xml.Bytes())
	if feedback.PolicyPublished.Domain == "" && feedback.ReportMetadata.ReportID == "" {
		response.Msg = "File not supported."
		json.NewEncoder(w).Encode(response)
		return

	}
	feedbackJson, _ := json.Marshal(feedback)

	feedbackJsonStr := string(feedbackJson)
	feedbackJsonStr = strings.ReplaceAll(feedbackJsonStr, "\n", "")
	fmt.Println(feedbackJsonStr)
	DB.InsertIntoDB(feedback.ReportMetadata.ReportID, feedback.PolicyPublished.Domain, feedbackJsonStr, userId)

	response.Status = true
	response.Msg = "Success"
	response.ReportID = feedback.ReportMetadata.ReportID

	json.NewEncoder(w).Encode(response)

}
