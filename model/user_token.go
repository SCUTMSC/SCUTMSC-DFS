package model

import (
	"fmt"

	"../model/mysql"
)

// UpdateUserToken is to update user token once trying to sign in
func UpdateUserToken(nickname string, token string) bool {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare(
		"REPLACE INTO tbl_user_token (`user_name`, `user_token`) VALUES (?, ?)")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when updating user token, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when updating user token, err: \n" + err.Error())
		return false
	}
	defer stmt.Close()

	// Execute SQL statement
	_, err = stmt.Exec(nickname, token)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when updating user token, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when updating user token, err: \n" + err.Error())
		return false
	}
	return true
}

// CheckUserToken is to check whether the user token is valid
func CheckUserToken(nickname string, token string) bool {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare("SELECT user_token FROM tbl_user_token WHERE user_name=?")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when checking user token, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when checking user token, err: \n" + err.Error())
		return false
	}
	defer stmt.Close()

	// Execute SQL statement
	var tokDB string
	err = stmt.QueryRow(nickname).Scan(&tokDB)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when checking user token, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when checking user token, err: \n" + err.Error())
		return false
	}

	// Compare results with parameters
	if len(tokDB) > 0 && tokDB == token {
		return true
	}
	return false
}
