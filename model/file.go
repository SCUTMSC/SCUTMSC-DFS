package model

import (
	"database/sql"
	"fmt"

	"../model/mysql"
)

// FileRecord is a database record of a file
type FileRecord struct {
	FileSha1    string
	FileName    sql.NullString
	FileSize    sql.NullInt64
	FilePath    sql.NullString
	EnableTimes sql.NullInt64
	EnableDays  sql.NullInt64
	CreateAt    sql.NullString
	UpdateAt    sql.NullString
}

// InsertFileRecord is to insert a database record of a file
func InsertFileRecord(fileSha1 string, fileName string, fileSize int64, filePath string,
	enableTimes int64, enableDays int64, createAt string, updateAt string) bool {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare(
		"INSERT ignore INTO tbl_file (`file_sha1`, `file_name`, `file_size`, `file_path`, " +
			"`enable_times`, `enable_days`, `create_at`, `update_at`, `status`) values (?, ?, ?, ?, ?, ?, ?, ?, 1)")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when inserting, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when inserting, err: \n" + err.Error())
		return false
	}
	defer stmt.Close()

	// Execute SQL statement
	res, err := stmt.Exec(fileSha1, fileName, fileSize, filePath, enableTimes, enableDays, createAt, updateAt)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when inserting, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when inserting, err: \n" + err.Error())
		return false
	}

	// Inform duplicate insert
	if rows, err := res.RowsAffected(); err == nil {
		// If rows == 0, duplicate insert occurs
		if rows == 0 {
			// log.Fatal("Duplicate insert occurs, file hash: \n" + fileSha1)
			fmt.Printf("Duplicate insert occurs, file hash: \n" + fileSha1)
		}
		// If rows == 1, duplicate insert didn't occur
		return true
	}
	return false
}

// UpdateFileRecord is to update a database record of a file
func UpdateFileRecord(fileSha1 string, fileName string, filePath string) bool {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare(
		"Update tbl_file SET file_name=?, file_path=? WHERE file_sha1=? AND status=1 LIMIT 1")
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when updating, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when updating, err: \n" + err.Error())
		return false
	}
	defer stmt.Close()

	// Execute SQL statement
	res, err := stmt.Exec(fileName, filePath, fileSha1)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when updating, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when updating, err: \n" + err.Error())
		return false
	}

	// Inform record not found
	if rows, err := res.RowsAffected(); err == nil {
		// If rows == 0, record not found
		if rows == 0 {
			// log.Fatal("Record not found, file hash: \n" + fileSha1)
			fmt.Printf("Record not found, file hash: \n" + fileSha1)
		}
		// If rows == 1, record found
		return true
	}
	return false
}

// CheckFileRecord is to check whether a database record of file exists
func CheckFileRecord(fileSha1 string) bool {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare("SELECT 1 FROM tbl_file where file_sha1=? AND status=1 LIMIT 1")
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when checking, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when checking, err: \n" + err.Error())
		return false
	}
	defer stmt.Close()

	// Execute SQL statement
	rows, err := stmt.Query(fileSha1)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when checking, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when checking, err: \n" + err.Error())
		return false
	}

	// Inform record not found
	if rows == nil || !rows.Next() {
		// log.Fatal("Record not found, file hash: \n" + fileSha1)
		fmt.Printf("Record not found, file hash: \n" + fileSha1)
		return false
	}
	return true
}

// SelectFileRecord is to select a database record of a file by fileSha1
func SelectFileRecord(fileSha1 string) (*FileRecord, bool) {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare(
		"SELECT file_sha1, file_name, file_size, file_path, enable_times, enable_days, create_at, update_at " +
			"FROM tbl_file WHERE file_sha1=? AND status=1 LIMIT 1")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when selecting, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when selecting, err: \n" + err.Error())
		return nil, false
	}
	defer stmt.Close()

	// Execute SQL statement
	fileRecord := FileRecord{}
	err = stmt.QueryRow(fileSha1).Scan(&fileRecord.FileSha1, &fileRecord.FileName, &fileRecord.FileSize, &fileRecord.FilePath,
		&fileRecord.EnableTimes, &fileRecord.EnableDays, &fileRecord.CreateAt, &fileRecord.UpdateAt)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when selecting, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when selecting, err: \n" + err.Error())
		return nil, false
	}

	// Return query result
	return &fileRecord, true
}

// SelectFileRecords is to select database records of files by limitCount
func SelectFileRecords(limitCount int) ([]FileRecord, bool) {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare(
		"SELECT file_sha1, file_name, file_size, file_path, enable_times, enable_days, create_at, update_at " +
			"FROM tbl_file WHERE status=1 LIMIT ?")
	if err != nil {
		// log.Fatal("Failed to prepare SQL statement when selecting, err: \n" + err.Error())
		fmt.Printf("Failed to prepare SQL statement when selecting, err: \n" + err.Error())
		return nil, false
	}
	defer stmt.Close()

	// Execute SQL statement
	rows, err := stmt.Query(limitCount)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when selecting, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when selecting, err: \n" + err.Error())
		return nil, false
	}
	var fileRecords []FileRecord
	for rows.Next() {
		fileRecord := FileRecord{}
		err = rows.Scan(&fileRecord.FileSha1, &fileRecord.FileName, &fileRecord.FileSize, &fileRecord.FilePath,
			&fileRecord.EnableTimes, &fileRecord.EnableDays, &fileRecord.CreateAt, &fileRecord.UpdateAt)
		if err != nil {
			// log.Fatal("Failed to execute SQL statement when selecting, err: \n" + err.Error())
			fmt.Printf("Failed to execute SQL statement when selecting, err: \n" + err.Error())
			return nil, false
		}
		fileRecords = append(fileRecords, fileRecord)
	}

	// Return query results
	return fileRecords, true
}

// DeleteFileRecord is to delete a database record of a file
func DeleteFileRecord(fileSha1 string) bool {
	// Prepare SQL statement
	stmt, err := mysql.DBConn().Prepare(
		"UPDATE tbl_file SET status=2 WHERE file_sha1=? AND status=1 LIMIT 1")
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when deleting, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when deleting, err: \n" + err.Error())
		return false
	}
	defer stmt.Close()

	// Execute SQL statement
	res, err := stmt.Exec(fileSha1)
	if err != nil {
		// log.Fatal("Failed to execute SQL statement when deleting, err: \n" + err.Error())
		fmt.Printf("Failed to execute SQL statement when deleting, err: \n" + err.Error())
		return false
	}

	// Inform record not found
	if rows, err := res.RowsAffected(); err == nil {
		// If rows == 0, record not found
		if rows == 0 {
			// log.Fatal("Record not found, file hash: \n" + fileSha1)
			fmt.Printf("Record not found, file hash: \n" + fileSha1)
		}
		// If rows == 1, record found
		return true
	}
	return false
}
