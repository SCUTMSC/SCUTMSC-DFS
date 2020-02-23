package meta

// MPUploadMeta is the metadata of a multipart upload try
type MPUploadMeta struct {
	FileSha1   string
	FileSize   int64
	UploadID   string
	ChunkSize  int
	ChunkCount int
}
