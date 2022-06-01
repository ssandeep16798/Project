package DB

import (
	"database/sql"
	"dmarc_backend/internal/domain/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "dmarc"
)

func GetDmarcFeedbackData(reportId string, domainName string) []models.Data {

	var data *sql.Rows

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	getDmarcFeedbackData := `SELECT id, dmarc_json FROM dmarc_schema.dmarc_data WHERE `
	countQuery := `SELECT COUNT(*) FROM dmarc_schema.dmarc_data WHERE `
	fmt.Println("reportId: " + reportId)
	fmt.Println("domainName: " + domainName)
	var count int

	if reportId != "" && domainName != "" {
		fmt.Println("Fetch by Domain and Report ID")

		getDmarcFeedbackData = getDmarcFeedbackData + `dmarc_json -> 'ReportMetadata' ->> 'ReportID'= $1` + ` AND ` + `dmarc_json -> 'PolicyPublished' ->> 'Domain' = $2`
		countQuery = countQuery + `dmarc_json -> 'ReportMetadata' ->> 'ReportID'= $1` + ` AND ` + `dmarc_json -> 'PolicyPublished' ->> 'Domain' = $2`
		err = db.QueryRow(countQuery, reportId, domainName).Scan(&count)
		data, err = db.Query(getDmarcFeedbackData, reportId, domainName)
		CheckError(err)
	} else if reportId == "" && domainName != "" {
		fmt.Println("Fetch Only by Domain")
		getDmarcFeedbackData = getDmarcFeedbackData + `dmarc_json -> 'PolicyPublished' ->> 'Domain' = $1`
		countQuery = countQuery + `dmarc_json -> 'PolicyPublished' ->> 'Domain' = $1`
		err = db.QueryRow(countQuery, domainName).Scan(&count)
		data, err = db.Query(getDmarcFeedbackData, domainName)
		CheckError(err)
	} else {
		fmt.Println("Fetch Only By Report ID")
		getDmarcFeedbackData = getDmarcFeedbackData + `dmarc_json -> 'ReportMetadata' ->> 'ReportID'= $1`
		countQuery = countQuery + `dmarc_json -> 'ReportMetadata' ->> 'ReportID'= $1`
		fmt.Println(countQuery)
		err = db.QueryRow(countQuery, reportId).Scan(&count)
		fmt.Println(count)
		data, err = db.Query(getDmarcFeedbackData, reportId)
		CheckError(err)
	}

	cols, err := data.Columns()
	CheckError(err)
	rowSize := len(cols)

	dmarcData := make([]models.Data, rowSize)
	defer data.Close()
	dmarcData = MapRows(data, count)

	return dmarcData
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func MapRows(rows *sql.Rows, count int) []models.Data {

	// Create our map, and retrieve the value for each column from the pointers slice,
	// storing it in the map with the name of the column as the key.

	cols, _ := rows.Columns()
	data := make([]models.Data, count)
	index := 0
	for rows.Next() {
		m := make(map[string]interface{})
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			CheckError(err)
		}

		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		str := string(m["dmarc_json"].([]uint8))
		id := m["id"].(int64)
		data[index].Id = id
		var dmarcData models.Feedback
		json.Unmarshal([]byte(str), &dmarcData)
		data[index].Dmarc_json = dmarcData
		fmt.Println(str)
		fmt.Println(id)
		index += 1

	}
	return data

}

func GetDmarcFeedbackDataByUserId(userId int) models.UserId {

	var data *sql.Rows
	var useridresponse models.UserId
	var count int

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	isAdmin := CheckIfAdminByUserId(userId)
	fmt.Println(isAdmin)
	if isAdmin {
		getDmarcFeedbackData := `SELECT dmarc_schema.users.id,firstname,lastname,username,dmarc_schema.users.created_datetime,report_id,domain_name FROM dmarc_schema.users  JOIN dmarc_schema.dmarc_data ON dmarc_schema.dmarc_data.user_id=dmarc_schema.users.id`
		countQuery := `SELECT COUNT(*) FROM dmarc_schema.users RIGHT JOIN dmarc_schema.dmarc_data ON dmarc_schema.dmarc_data.user_id=dmarc_schema.users.id`

		fmt.Println("Fetch All by Admin")

		fmt.Println(getDmarcFeedbackData)
		err = db.QueryRow(countQuery).Scan(&count)
		fmt.Println(count)
		if count < 1 {
			useridresponse.Status = false
			useridresponse.Data = make([]models.Users, 0)
			useridresponse.Msg = ""
			return useridresponse
		}
		fmt.Println("After count query")
		fmt.Println("COuntquery%v", countQuery)
		data, err = db.Query(getDmarcFeedbackData)

	} else {

		getDmarcFeedbackData := `SELECT dmarc_schema.users.id,firstname,lastname,username,dmarc_schema.users.created_datetime,report_id,domain_name FROM dmarc_schema.users, dmarc_schema.dmarc_data WHERE dmarc_schema.users.id=$1 AND dmarc_schema.dmarc_data.user_Id=$1`
		countQuery := `SELECT COUNT(*) FROM dmarc_schema.users,dmarc_schema.dmarc_data WHERE dmarc_schema.users.id=$1 AND dmarc_schema.dmarc_data.user_Id=$1`
		fmt.Println("user_id: " + strconv.Itoa(userId))

		fmt.Println("Fetch by User ID")

		fmt.Println(getDmarcFeedbackData)
		err = db.QueryRow(countQuery, userId).Scan(&count)
		fmt.Println(count)
		if count < 1 {
			useridresponse.Status = false
			useridresponse.Data = make([]models.Users, 0)
			useridresponse.Msg = ""
			return useridresponse
		}
		fmt.Println("After count query")
		fmt.Println("COuntquery%v", countQuery)
		data, err = db.Query(getDmarcFeedbackData, userId)
	}
	CheckError(err)
	fmt.Println("After Data query")
	cols, err := data.Columns()
	CheckError(err)
	rowSize := len(cols)

	usersData := make([]models.Users, rowSize)
	defer data.Close()
	usersData = MapUserData(data, count)
	fmt.Println("After mapping")
	useridresponse.Status = true
	useridresponse.Msg = "Success"
	useridresponse.Data = usersData

	return useridresponse
}

func MapUserData(rows *sql.Rows, count int) []models.Users {

	// Create our map, and retrieve the value for each column from the pointers slice,
	// storing it in the map with the name of the column as the key.

	cols, _ := rows.Columns()
	fmt.Printf("count %v:\n", count)
	data := make([]models.Users, count)
	index := 0
	for rows.Next() {
		fmt.Printf("In iteration %v\n", index)
		fmt.Println(data)
		m := make(map[string]interface{})
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		fmt.Println(len(cols))
		columnPointers := make([]interface{}, len(cols))

		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			CheckError(err)
		}

		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		fmt.Printf("After scanning coloumns %v\n", index)

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		report_id := (m["report_id"].(string))
		id := m["id"].(int64)
		firstname := (m["firstname"].(string))
		lastname := (m["lastname"].(string))
		username := firstname + " " + lastname
		domain_name := (m["domain_name"].(string))
		created_datetime := (m["created_datetime"].(time.Time))

		fmt.Println(len(data))
		data[index].Id = id
		data[index].Report_id = report_id
		data[index].Username = username
		data[index].Domain_name = domain_name
		data[index].User_id = id
		data[index].Created_datetime = created_datetime

		fmt.Println(id)
		index += 1

	}
	return data
}

func CheckIfAdminByUserId(user_id int) bool {

	var userType int
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()
	fetchDynamicStatement := `SELECT "user_type" FROM dmarc_schema.users WHERE id=$1`
	err = db.QueryRow(fetchDynamicStatement, user_id).Scan(&userType)
	CheckError(err)

	fmt.Println(userType)
	return userType == 0

}
