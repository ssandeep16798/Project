package structconverter

import (
	"dmarc_backend/internal/domain/models"
	"encoding/xml"
	"fmt"
)

//Convert the xml file to struct with xml file path as a parameter
func ConvertToStruct(byteValue []byte) models.Feedback {

	var feedback models.Feedback
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &feedback)

	fmt.Println(feedback)

	return feedback
}
