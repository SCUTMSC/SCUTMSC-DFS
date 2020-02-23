package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	db "../model"
	"github.com/julienschmidt/httprouter"
)

// IndexHandler is to handle directing to index.html
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := ioutil.ReadFile("./static/view/signup.html")
	if err != nil {
		log.Fatal("Failed to read index.html, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileUploadHandler is to handle browser clients uploading files to the http server.
func FileUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	fileSha1 := r.Form.Get("fileSha1")
	fileSize, _ := strconv.Atoi(r.Form.Get("fileSize"))

	// Check whether the file has already existed
	if isExist := db.CheckFileRecord(fileSha1); isExist {
		// If file exists, try fast upload method
		fmt.Println("Switch to fast upload mode...")
		FileFastUploadHandler(w, r, ps)
	} else {
		// If file doesn't exist, check whether the file is large enough
		if fileSize > 5*1024*1024 {
			// If file is larger than 5MB, try multipart upload method
			fmt.Println("Switch to multipart upload mode...")
			FileMPUploadInitHandler(w, r, ps)
		} else {
			// If file isn't larger than 5MB, try normal upload method
			fmt.Println("Switch to normal upload mode...")
			FileNormalUploadHandler(w, r, ps)
		}
	}
}

// FileDownloadHandler is to handle browser clients downloading files from the http server.
func FileDownloadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	fileSize, _ := strconv.Atoi(r.Form.Get("fileSize"))

	// Check whether the file is large enough
	if fileSize > 5*1024*1024 {
		// If file is larger than 5MB, try breakpoint-resumed download method
		fmt.Println("Switch to breakpoint-resumed download mode...")
		FileBRDownloadHandler(w, r, ps)
	} else {
		// If file isn't larger than 5MB, try normal download method
		fmt.Println("Switch to normal download mode...")
		FileDownloadAttachmentHandler(w, r, ps)
	}
}
