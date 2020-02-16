package meta

import (
	"sort"
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

// SetFileMeta is to update the metadata of a file
func SetFileMeta(fileSha1 string, fileMeta FileMeta) {
	key := fileSha1
	value := fileMeta
	fileMetas[key] = value
}

// GetFileMeta is to retrieve the metadata of a file by fileSha1
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
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

// DeleteFileMeta is to update the metadata of a file
func DeleteFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
