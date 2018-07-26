package user

import (
	"fmt"
	"log"

	"github.com/mrsmuneton/platform-test/src/db"
	"github.com/mrsmuneton/platform-test/src/utils"
	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	Code string `code`
}

//using timestamsp ca improve sorting efficiency in queries
type User struct {
	Id              int    `id`
	CreatedDate     string `createdDate` //cheating by a string, this must be a timestamp
	CurrentPassword string `currentPassword`
	Email           string `email`
	Name            string `name`
	UpdatedDate     string `updatedDate` //cheating by a string, this must be a timestamp
}

func BcryptEncrypt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func BcryptMatchPassword(storedHash string, enteredPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(enteredPassword)); err != nil {
		return false
	}
	return true
}

func CreateUser(newUser User) bool {
	var error_result = false

	// _, e := ValidateEmailIsUnique(newUser)
	// if e.Code != "" {
	// 	return true
	// }

	_, e := ValidateUserMinimumFields(newUser)
	if e.Code != "" {
		return true
	}

	dbConnection, err := db.DBConnect()
	if err != nil {
		fmt.Println(err)
		error_result = true
	}

	password, err := BcryptEncrypt(newUser.CurrentPassword) // TODO: This should implement base64 from the client request
	t := utils.CurrentTime()

	_, queryerr := dbConnection.Query("INSERT INTO users(created_date, currentpassword, email, name, updated_date) VALUES($1,$2,$3,$4,$5);", t, password, newUser.Email, newUser.Name, t)
	if queryerr != nil {
		fmt.Println(queryerr)
		error_result = true
	}

	return error_result
}

func DeleteUser(id string) bool {
	var error_result = false

	dbConnection, err := db.DBConnect()
	if err != nil {
		fmt.Println(err)
		error_result = true
	}
	query := "DELETE FROM users WHERE id=$1"
	_, queryerr := dbConnection.Query(query, id)
	if queryerr != nil {
		fmt.Println(queryerr)
		error_result = true
	}

	return error_result
}

func GetUserRecordById(user_id string) (User, bool) {
	var error_result bool
	var userRecord User
	dbConnection, err := db.DBConnect()
	if err != nil {
		fmt.Println(err)
		error_result = true
	}
	var query = "SELECT id, created_date, email, name, updated_date FROM users WHERE id=$1"
	queryerr := dbConnection.QueryRow(query, user_id).Scan(&userRecord.Id, &userRecord.CreatedDate, &userRecord.Email, &userRecord.Name, &userRecord.UpdatedDate)
	if queryerr != nil {
		fmt.Println(queryerr)
		error_result = true
	}
	return userRecord, error_result
}

func GetUserRecordByEmail(email string) (User, bool) {
	var error_result bool
	var userRecord User
	dbConnection, err := db.DBConnect()
	if err != nil {
		fmt.Println("DB connection issue", err)
		error_result = true
	}
	var query = "SELECT id, created_date, email, name, updated_date FROM users WHERE email=$1"
	queryerr := dbConnection.QueryRow(query, email).Scan(&userRecord.Id, &userRecord.CreatedDate, &userRecord.Email, &userRecord.Name, &userRecord.UpdatedDate)
	if queryerr != nil {
		fmt.Println(queryerr)
		error_result = true
	}
	return userRecord, error_result
}

func LoginUser(userRequest User) (User, bool) {
	dbConnection, queryerr := db.DBConnect()
	var userRecord User
	var query = "SELECT id, currentpassword, email, name FROM users WHERE email=$1"
	err := dbConnection.QueryRow(query, userRequest.Email).Scan(&userRecord.Id, &userRecord.CurrentPassword, &userRecord.Email, &userRecord.Name)
	if queryerr != nil || err != nil {
		fmt.Println("LoginUser Query Error")
		return User{}, false
	}
	if BcryptMatchPassword(userRecord.CurrentPassword, userRequest.CurrentPassword) {
		return userRecord, true
	} else {
		return User{}, false
	}
}

func UpdateUserFields(user_id string, u User) (User, bool) {
	var error_result = false
	_, e := ValidateUserMinimumFields(u)
	if e.Code != "" {
		return u, true
	}

	dbConnection, err := db.DBConnect()
	if err != nil {
		fmt.Println(err)
		error_result = true
	}
	var now = utils.CurrentTime()
	var query = "UPDATE users SET currentpassword=$1, email=$2, name=$3, updated_date=$4 WHERE id=$5 RETURNING id, created_date, email, name, updated_date"
	queryerr := dbConnection.QueryRow(query, u.CurrentPassword, u.Email, u.Name, now, string(user_id)).Scan(&u.Id, &u.CreatedDate, &u.Email, &u.Name, &u.UpdatedDate)
	if queryerr != nil {
		fmt.Println(queryerr)
		error_result = true
	}
	return u, error_result
}

func ValidateUserMinimumFields(u User) (User, Error) {
	var requiredFields string
	e := Error{Code: ""}
	fmt.Print(e)

	if u.CurrentPassword == "" {
		requiredFields = requiredFields + string(" CurrentPassword")
	}

	if u.Email == "" {
		requiredFields = requiredFields + string(" Email")
	}

	if u.Name == "" {
		requiredFields = requiredFields + string(" Name")
	}

	if len(requiredFields) > 0 {
		e.Code = "Please provide all fields, including:" + requiredFields
	}

	return u, e
}

func ValidateEmailIsUnique(email string) bool {
	var count int
	dbConnection, err := db.DBConnect()
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = dbConnection.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", email).Scan(&count)
	if err != nil {
		log.Fatal(err)
		return false
	}

	if count > 0 {
		return false
	}
	return true
}
