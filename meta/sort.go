package meta

import "time"

const baseFormat = "2006-01-02 15:04:05"

type ByUploadAt []FileMeta

func (a ByUploadAt) Len() int {
	return len(a)
}

func (a ByUploadAt) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByUploadAt) Less(i, j int) bool {
	iTime, _ := time.Parse(baseFormat, a[i].UploadAt)
	jTime, _ := time.Parse(baseFormat, a[j].UploadAt)
	return iTime.UnixNano() > jTime.UnixNano()
}
