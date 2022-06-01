package DB

import (
	"database/sql"
	"dmarc_backend/internal/domain/models"
	"fmt"

	_ "github.com/lib/pq"
)

//Checks for username and password in the Database.
func UserDB(username string, pass string) int {

	var count int
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	fetchuser := `SELECT COUNT(*) FROM dmarc_schema.users WHERE username=$1 AND password=$2`
	err = db.QueryRow(fetchuser, username, pass).Scan(&count)
	CheckError(err)

	fmt.Println(count)
	return count

}

//Checks if the email id entered by the user is present in the Database.
func CheckIfEmailExists(email string, username string) int {

	var count int
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	checkEmailQuery := `SELECT COUNT(*) FROM dmarc_schema.users WHERE email=$1 OR username=$2`
	err = db.QueryRow(checkEmailQuery, email, username).Scan(&count)
	CheckError(err)

	fmt.Println(count)
	return count

}

//New User Registration
func ResgisterUser(req models.UserRegistrationRequest) bool {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()
	insertDynStmt := `INSERT INTO dmarc_schema.users("firstname", "lastname", "email", "gender", "phone_number", "username", "password", "created_datetime" ) VALUES($1, $2, $3, $4, $5, $6, $7, NOW())`
	_, err = db.Exec(insertDynStmt, req.Firstname, req.Lastname, req.Email, req.Gender, req.Mobile, req.Username, req.Password)
	CheckError(err)

	if err != nil {
		return false
	} else {
		return true
	}
}

func FetchData(username string, pass string) models.UserRegistrationRequest {

	var registerReq models.UserRegistrationRequest

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()
	fetchDynamicStatement := `SELECT "firstname", "lastname", "email", "username", "password" , "id", "gender", "phone_number" FROM dmarc_schema.users WHERE username=$1 AND password=$2`

	err = db.QueryRow(fetchDynamicStatement, username, pass).Scan(&registerReq.Firstname, &registerReq.Lastname, &registerReq.Email, &registerReq.Username, &registerReq.Password, &registerReq.Id, &registerReq.Gender, &registerReq.Mobile)
	CheckError(err)

	return registerReq

}

func FetchAllData(username string, pass string) []models.UserRegistrationRequest {

	var count int

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()
	fetchDynamicStatement := `SELECT "firstname", "lastname", "email", "username", "password" , "id", "gender", "phone_number" FROM dmarc_schema.users`

	data, err := db.Query(fetchDynamicStatement)
	CheckError(err)

	cols, err := data.Columns()
	CheckError(err)
	rowSize := len(cols)

	countQuery := `SELECT COUNT(*) FROM dmarc_schema.users`
	fmt.Println(countQuery)
	err = db.QueryRow(countQuery).Scan(&count)
	fmt.Println(count)

	dmarcData := make([]models.UserRegistrationRequest, rowSize)
	defer data.Close()
	dmarcData = MapUserRegistrationData(data, count)

	return dmarcData
}

func EmailVerification(email string) int {

	var count int
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	fmt.Println(email)

	// dynamic
	//insertDynStmt := `INSERT INTO dmarc_schema.users("user_name", "password") VALUES($1, $2)`
	checkEmailQuery := `SELECT COUNT(*) FROM dmarc_schema.users WHERE email=$1`
	err = db.QueryRow(checkEmailQuery, email).Scan(&count)
	CheckError(err)

	fmt.Println(count)
	return count

}

func UpdatePassword(email string, pass string) bool {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	updatePasswordQuery := `UPDATE dmarc_schema.users SET password=$2 WHERE email=$1`
	_, err = db.Exec(updatePasswordQuery, email, pass)
	CheckError(err)

	if err != nil {
		return false
	} else {
		return true
	}

}

func RegisterData(email string) int64 {

	var userId int64
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	data := `SELECT id FROM dmarc_schema.dmarc_users WHERE email=$1`
	err = db.QueryRow(data, email).Scan(&data)
	CheckError(err)

	fmt.Println(data)
	return userId

}

func CheckIfAdmin(username string, pass string) bool {

	var userType int
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()
	fetchDynamicStatement := `SELECT "user_type" FROM dmarc_schema.users WHERE username=$1 AND password=$2`
	err = db.QueryRow(fetchDynamicStatement, username, pass).Scan(&userType)
	CheckError(err)

	fmt.Println(userType)
	return userType == 0

}

func MapUserRegistrationData(rows *sql.Rows, count int) []models.UserRegistrationRequest {

	// Create our map, and retrieve the value for each column from the pointers slice,
	// storing it in the map with the name of the column as the key.

	cols, _ := rows.Columns()
	data := make([]models.UserRegistrationRequest, count)
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
		//SELECT "firstname", "lastname", "email", "username", "password" , "id", "gender", "phone_number" FROM dmarc_schema.users`

		id := int(m["id"].(int64))
		firstname := (m["firstname"].(string))
		lastname := (m["lastname"].(string))
		email := (m["email"].(string))
		username := (m["username"].(string))
		password := (m["password"].(string))
		gender := (m["gender"].(string))
		phone_number := (m["phone_number"].(string))

		data[index].Id = id
		data[index].Username = username
		data[index].Firstname = firstname
		data[index].Lastname = lastname
		data[index].Email = email
		data[index].Gender = gender
		data[index].Password = password
		data[index].Mobile = phone_number

		fmt.Println(id)
		index += 1

	}
	return data
}
