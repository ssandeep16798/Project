package DB

import (
	"database/sql"
	"dmarc_backend/internal/domain/models"
	"fmt"
	"time"
)

func GetDmarcFeedBackDataByDomain(domain string, startdate string, enddate string) models.DomainResponse {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var domainResponse models.DomainResponse

	domainResponse.Data = make([]models.DomainData, 0)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	var count int
	var data *sql.Rows
	defer db.Close()
	fmt.Println("Fetch Only by Domain")
	getDmarcFeedbackData := `SELECT dmarc_json -> 'ReportMetadata' ->> 'ReportID' AS report_id, created_datetime FROM dmarc_schema.dmarc_data WHERE `
	countQuery := `SELECT COUNT(*) FROM dmarc_schema.dmarc_data WHERE `
	getDmarcFeedbackData = getDmarcFeedbackData + `dmarc_json -> 'PolicyPublished' ->> 'Domain' = $1 and date(created_datetime) BETWEEN $2 AND $3`
	countQuery = countQuery + `dmarc_json -> 'PolicyPublished' ->> 'Domain' = $1 and date(created_datetime) BETWEEN $2 AND $3`
	err = db.QueryRow(countQuery, domain, startdate, enddate).Scan(&count)

	if count < 1 {

		domainResponse.Status = false
		domainResponse.Msg = "No Records Found"

		return domainResponse
	}

	data, err = db.Query(getDmarcFeedbackData, domain, startdate, enddate)
	CheckError(err)

	fmt.Println(countQuery, getDmarcFeedbackData)

	defer data.Close()
	domainResponse.Data = MapDomainData(data, count)
	domainResponse.Status = true
	domainResponse.Msg = "Data Found"

	return domainResponse
}

func MapDomainData(rows *sql.Rows, count int) []models.DomainData {

	// Create our map, and retrieve the value for each column from the pointers slice,
	// storing it in the map with the name of the column as the key.

	cols, _ := rows.Columns()
	data := make([]models.DomainData, count)
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
		report_id := string(m["report_id"].(string))
		created_datetime := m["created_datetime"].(time.Time)
		data[index].ReportID = report_id

		data[index].Created_datetime = created_datetime
		fmt.Println(report_id)
		fmt.Println(created_datetime)
		index += 1

	}
	return data

}
