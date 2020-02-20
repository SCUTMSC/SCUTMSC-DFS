package model

import (
	"fmt"

	"../model/mysql"
)

// UserFileRecord is a database record of a user file
type UserFileRecord struct {
	Nickname string
	FileSha1 string
}

// AppendUserFile is to append user file to database
func AppendUserFile(nickname string, fileSha1 string) bool {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare(
		"INSERT ignore INTO tbl_user_file (`user_name`, `file_sha1`, `status`) VALUES (?, ?, 1)")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when appending user file, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when appending user file, err: \n" + err.Error())
		return false
	}
	defer stmt.Close()

	// Execute SQL statement
	_, err = stmt.Exec(nickname, fileSha1)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when appending user file, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when appending user file, err: \n" + err.Error())
		return false
	}
	return true
}

// GetUserFiles is to retrieve user files from database
func GetUserFiles(nickname string, limitCount int) ([]UserFileRecord, bool) {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare("SELECT user_name, file_sha1 FROM tbl_user_file WHERE user_name=? AND status=1 LIMIT ?")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when querying user files, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when querying user files, err: \n" + err.Error())
		return nil, false
	}
	defer stmt.Close()

	// Execute SQL statement
	rows, err := stmt.Query(nickname, limitCount)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when querying user files, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when querying user files, err: \n" + err.Error())
		return nil, false
	}
	var userFileRecords []UserFileRecord
	for rows.Next() {
		userFileRecord := UserFileRecord{}
		err = rows.Scan(&userFileRecord.Nickname, &userFileRecord.FileSha1)
		if err != nil {
			// log.Fatal("Failed to execute SQL statement when querying user files, err: \n" + err.Error())
			fmt.Printf("Failed to execute SQL statement when querying user files, err: \n" + err.Error())
			return nil, false
		}
		userFileRecords = append(userFileRecords, userFileRecord)
	}

	// Return query results
	return userFileRecords, true
}

// DeleteUserFile is to delete user file in database
func DeleteUserFile(nickname string, fileSha1 string) bool {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare(
		"UPDATE tbl_user_file SET status=2 WHERE user_name=? AND file_sha1=? LIMIT 1")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when deleting user file, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when deleting user file, err: \n" + err.Error())
		return false
	}
	defer stmt.Close()

	// Execute SQL statement
	res, err := stmt.Exec(nickname, fileSha1)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when deleting user file, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when deleting user file, err: \n" + err.Error())
		return false
	}

	// Inform record not found
	if rows, err := res.RowsAffected(); err == nil {
		// If rows == 0, record not found
		if rows == 0 {
			// log.Fatal("Record not found, nickname: \n" + nickname)
			// log.Fatal("Record not found, file hash: \n" + fileSha1)
			fmt.Printf("Record not found, nickname: \n" + nickname)
			fmt.Printf("Record not found, file hash: \n" + fileSha1)
		}
		// If rows == 1, record found
		return true
	}
	return false
}
