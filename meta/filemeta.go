package meta

import (
	"sort"

	db "../model"
)

// FileMeta is the metadata of a file
type FileMeta struct {
	FileSha1    string
	FileName    string
	FileSize    int64
	FilePath    string
	EnableTimes int64
	EnableDays  int64
	CreateAt    string
	UpdateAt    string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// CreateFileMeta is to create the metadata of a file
func CreateFileMeta(fileMeta FileMeta) {
	key := fileMeta.FileSha1
	value := fileMeta
	fileMetas[key] = value
}

// CreateFileMetaDB is to create the metadata of a file in database
func CreateFileMetaDB(fileMeta FileMeta) bool {
	return db.InsertFileRecord(fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize, fileMeta.FilePath,
		fileMeta.EnableTimes, fileMeta.EnableDays, fileMeta.CreateAt, fileMeta.UpdateAt)
}

// SetFileMeta is to update the metadata of a file
func SetFileMeta(fileSha1 string, fileMeta FileMeta) {
	key := fileSha1
	value := fileMeta
	fileMetas[key] = value
}

// SetFileMetaDB is to update the metadata of a file in database
func SetFileMetaDB(fileSha1 string, fileMeta FileMeta) bool {
	return db.UpdateFileRecord(fileSha1, fileMeta.FileName, fileMeta.FilePath)
}

// GetFileMeta is to retrieve the metadata of a file by fileSha1
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB is to retrieve the metadata of a file by fileSha1 in database
func GetFileMetaDB(fileSha1 string) (FileMeta, bool) {
	fileRecord, ok := db.SelectFileRecord(fileSha1)
	if !ok {
		return FileMeta{}, false
	}

	fileMeta := FileMeta{
		FileSha1:    fileRecord.FileSha1,
		FileName:    fileRecord.FileName.String,
		FileSize:    fileRecord.FileSize.Int64,
		FilePath:    fileRecord.FilePath.String,
		EnableTimes: fileRecord.EnableTimes.Int64,
		EnableDays:  fileRecord.EnableDays.Int64,
		CreateAt:    fileRecord.CreateAt.String,
		UpdateAt:    fileRecord.UpdateAt.String,
	}
	return fileMeta, true
}

// GetFileMetasByCreateAt is to retrieve the metadata of files by CreateAt
func GetFileMetasByCreateAt(limitCount int) []FileMeta {
	// Convert map to array of file metadata
	var fileMetaArray []FileMeta
	for _, fileMeta := range fileMetas {
		fileMetaArray = append(fileMetaArray, fileMeta)
	}

	// Sort by upload time
	sort.Sort(ByCreateAt(fileMetaArray))

	// Return the specified number of files
	return fileMetaArray[0:limitCount]
}

// GetFileMetasByCreateAtDB is to retrieve the metadata of files by CreateAt in database
func GetFileMetasByCreateAtDB(limitCount int) ([]FileMeta, bool) {
	fileRecords, ok := db.SelectFileRecords(limitCount)
	if !ok {
		return make([]FileMeta, 0), false
	}

	fileMetaArray := make([]FileMeta, len(fileRecords))
	for i := 0; i < len(fileRecords); i++ {
		fileMetaArray[i] = FileMeta{
			FileSha1:    fileRecords[i].FileSha1,
			FileName:    fileRecords[i].FileName.String,
			FileSize:    fileRecords[i].FileSize.Int64,
			FilePath:    fileRecords[i].FilePath.String,
			EnableTimes: fileRecords[i].EnableTimes.Int64,
			EnableDays:  fileRecords[i].EnableDays.Int64,
			CreateAt:    fileRecords[i].CreateAt.String,
			UpdateAt:    fileRecords[i].UpdateAt.String,
		}
	}
	return fileMetaArray, true
}

// DeleteFileMeta is to delete the metadata of a file
func DeleteFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

// DeleteFileMetaDB is to delete the metadata of a file in database
func DeleteFileMetaDB(fileSha1 string) bool {
	return db.DeleteFileRecord(fileSha1)
}
