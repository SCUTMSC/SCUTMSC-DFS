package meta

import "time"

const baseFormat = "2006-01-02 15:04:05"

type ByCreateAt []FileMeta

func (a ByCreateAt) Len() int {
	return len(a)
}

func (a ByCreateAt) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByCreateAt) Less(i, j int) bool {
	iTime, _ := time.Parse(baseFormat, a[i].CreateAt)
	jTime, _ := time.Parse(baseFormat, a[j].CreateAt)
	return iTime.UnixNano() > jTime.UnixNano()
}
