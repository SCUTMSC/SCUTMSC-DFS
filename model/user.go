package model

import (
	"fmt"

	"../model/mysql"
)

// UserRecord is a database record of a user
type UserRecord struct {
	Nickname string

	Realname   string
	Gender     int
	Birthday   string
	Campus     string
	School     string
	Major      string
	Grade      string
	Class      string
	Dormitory  string
	Department string

	Email  string
	Phone  string
	Wechat string
	QQ     string

	SignUpAt string
	SignInAt string

	Status int
}

// UserSignUp is to insert a database record of a user
func UserSignUp(nickname string, password string) bool {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare(
		"INSERT ignore INTO tbl_user (`nick_name`, `user_pwd`) VALUES (?, ?)")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when signuping, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when signuping, err: \n" + err.Error())
		return false
	}
	defer stmt.Close()

	// Execute SQL statement
	res, err := stmt.Exec(nickname, password)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when signuping, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when signuping, err: \n" + err.Error())
		return false
	}

	// Inform duplicate insert
	if rows, err := res.RowsAffected(); err == nil {
		// If rows == 0, duplicate insert occurs
		if rows == 0 {
			// log.Fatal("Duplicate insert occurs, nickname: \n" + nickname)
			fmt.Printf("Duplicate insert occurs, nickname: \n" + nickname)
		}
		// If rows == 1, duplicate insert didn't occur
		return true
	}
	return false
}

// UserSignIn is to select a database record of a user and compare it with params
func UserSignIn(nickname string, password string) bool {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare("SELECT user_pwd FROM tbl_user where nick_name=? LIMIT 1")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when signining, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when signining, err: \n" + err.Error())
		return false
	}
	defer stmt.Close()

	// Execute SQL statement
	var pwdDB string
	err = stmt.QueryRow(nickname).Scan(&pwdDB)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when signining, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when signining, err: \n" + err.Error())
		return false
	}

	// Compare results with parameters
	if len(pwdDB) > 0 && pwdDB == password {
		return true
	}
	return false
}

// GetUserInfo is to retrieve user information from the database
func GetUserInfo(nickname string) (*UserRecord, bool) {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare(
		"SELECT nick_name, signup_at FROM tbl_user WHERE nick_name=? LIMIT 1")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when querying, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when querying, err: \n" + err.Error())
		return nil, false
	}
	defer stmt.Close()

	// Execute SQL statement
	userRecord := UserRecord{}
	err = stmt.QueryRow(nickname).Scan(&userRecord.Nickname, &userRecord.SignUpAt)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when querying, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when querying, err: \n" + err.Error())
		return nil, false
	}

	// Return query result
	return &userRecord, true
}
