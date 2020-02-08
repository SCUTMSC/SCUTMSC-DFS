package controller

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
		FileName: fileHeader.Filename,
		FilePath: "./storage/tmp/" + fileHeader.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Create a local file to store the uploaded file
	err = os.MkdirAll("./storage/tmp/", os.ModePerm)
	localFile, err := os.Create(fileMeta.FilePath)
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

	// Return the http response to show the uploaded file
	data, err := json.Marshal(fileMeta)
	if err != nil {
		log.Fatal("Failed to convert fileMeta to JSON, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileUpdate is to handle updating files' metadata.
func FileUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	optionType := r.Form.Get("optionType")
	fileSha1 := r.Form.Get("fileHash")
	fileName := r.Form.Get("fileName")
	fileMeta := meta.GetFileMeta(fileSha1)

	// Check whether the option type is legal
	if optionType != "0" {
		log.Fatal("Failed to access rights to rename files, err: \n" + "option type error occurs")

		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Update the metadata of a file
	fileMeta.FileName = fileName
	meta.SetFileMeta(fileSha1, fileMeta)

	// Return the http response
	data, err := json.Marshal(fileMeta)
	if err != nil {
		log.Fatal("Failed to convert fileMeta to JSON, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileDownload is to handle browser clients downloading files from the http server.
func FileDownload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	fileSha1 := ps.ByName("fileHash")
	fileMeta := meta.GetFileMeta(fileSha1)

	// Open the local file to prepare the downloading file
	file, err := os.Open(fileMeta.FilePath)
	if err != nil {
		log.Fatal("Failed to open file when downloading, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the file from physical disks
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Failed to read file when downloading, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the http response
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+"\""+fileMeta.FileName+"\"")
	w.Write(data)
}

// SingleFileQuery is to handle querying files' metadata by fileHash.
func SingleFileQuery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	fileSha1 := ps.ByName("fileHash")

	// Use file hash to query file's metadata
	fileMeta := meta.GetFileMeta(fileSha1)
	data, err := json.Marshal(fileMeta)
	if err != nil {
		log.Fatal("Failed to convert fileMeta to JSON, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the http response
	w.Write(data)
}

// BatchFilesQuery is to handle querying files' metadata by limitCount.
func BatchFilesQuery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	limitCount, _ := strconv.Atoi(ps.ByName("limitCount"))

	// Use limit count to query files' metadata
	fileMetas := meta.GetFileMetasByUploadAt(limitCount)
	data, err := json.Marshal(fileMetas)
	if err != nil {
		log.Fatal("Failed to convert fileMetas to JSON, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the http response
	w.Write(data)
}

// FileDelete is to handle browser clients deleting files on the http server.
func FileDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ...
}
