package DB

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

func InsertIntoDB(reportId string, domainName string, feedbackJson string, userId int) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	feedbackJson = strings.ReplaceAll(feedbackJson, "\n", "")

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	// dynamic
	insertDynStmt := `INSERT INTO dmarc_schema.dmarc_data("report_id", "domain_name", "dmarc_json", "user_id", "created_datetime") VALUES($1, $2, $3, $4, NOW())`
	_, e := db.Exec(insertDynStmt, reportId, domainName, feedbackJson, userId)
	CheckError(e)
}
