package controller

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	"../meta"
	"../util"
)

// Index is to handle directing to upload.html
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := ioutil.ReadFile("./view/upload.html")
	if err != nil {
		log.Fatal("Failed to read upload.html, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileUpload is to handle browser clients uploading files to the http server.
func FileUpload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request to get the uploaded file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Fatal("Failed to read file when uploading, err: \n" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Define the metadata of a file
	fileMeta := meta.FileMeta{
		FileName:     fileHeader.Filename,
		FileLocation: "./storage/tmp/" + fileHeader.Filename,
		UploadAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	// Create a local file to store the uploaded file
	err = os.MkdirAll("./storage/tmp/", os.ModePerm)
	localFile, err := os.Create(fileMeta.FileLocation)
	if err != nil {
		log.Fatal("Failed to create file when uploading, err: \n" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	// Write the file into physical disks
	fileMeta.FileSize, err = io.Copy(localFile, file)
	if err != nil {
		log.Fatal("Failed to write file when uploading, err: \n" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Save the metadata of a file
	localFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(localFile)
	meta.CreateFileMeta(fileMeta)

	// Return the http response to inform the uploading success
	w.Write([]byte("Uploading files succeeded."))
}

// FileUpdate is to handle updating files' metadata.
func FileUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ...
}

// FileDownload is to handle browser clients downloading files from the http server.
func FileDownload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ...
}

// FileQuery is to handle querying files' metadata.
func FileQuery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ...
}

// FileDelete is to handle browser clients deleting files on the http server.
func FileDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ...
}
