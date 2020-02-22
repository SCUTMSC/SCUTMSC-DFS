package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"../meta"
	db "../model"
	"../util"
)

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
			// TODO: FileMPUploadInitHanlder(w, r, ps)
		} else {
			// If file isn't larger than 5MB, try normal upload method
			fmt.Println("Switch to normal upload mode...")
			FileNormalUploadHandler(w, r, ps)
		}
	}
}

// FileFastUploadHandler is to handle browser clients uploading files to the http server in fast mode.
func FileFastUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	nickname := r.Form.Get("nickname")
	fileSha1 := r.Form.Get("fileSha1")

	// Update user file info
	ok := db.AppendUserFile(nickname, fileSha1)

	// Return the http response
	if ok {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "SUCCESS",
		}
		w.Write(resp.JSONBytes())
	} else {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "FAILED",
		}
		w.Write(resp.JSONBytes())
	}
}

// FileNormalUploadHandler is to handle browser clients uploading files to the http server in normal mode.
func FileNormalUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request to get the uploaded file
	r.ParseForm()

	nickname := r.Form.Get("nickname")

	var enableTimes int
	var enableDays int
	if r.Form.Get("enable_times") == "" {
		enableTimes = 9999999 // Enable to download files 9999999 times in default
	} else {
		enableTimes, _ = strconv.Atoi(r.Form.Get("enable_times"))
	}
	if r.Form.Get("enable_days") == "" {
		enableDays = 30 // Enable to download files in 30 days in default
	} else {
		enableDays, _ = strconv.Atoi(r.Form.Get("enable_days"))
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Fatal("Failed to read file when uploading, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Define the metadata of a file
	fileMeta := meta.FileMeta{
		FileName:    fileHeader.Filename,
		FilePath:    "./storage/tmp/" + fileHeader.Filename,
		EnableTimes: int64(enableTimes),
		EnableDays:  int64(enableDays),
		CreateAt:    time.Now().Format("2006-01-02 15:04:05"),
		UpdateAt:    time.Now().Format("2006-01-02 15:04:05"),
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
	// meta.CreateFileMeta(fileMeta)
	_ = meta.CreateFileMetaDB(fileMeta)
	_ = db.AppendUserFile(nickname, fileMeta.FileSha1)

	// Return the http response to show the uploaded file
	data, err := json.Marshal(fileMeta)
	if err != nil {
		log.Fatal("Failed to convert fileMeta to JSON, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileUpdateHandler is to handle updating files' metadata.
func FileUpdateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	optionType := r.Form.Get("optionType")
	fileSha1 := r.Form.Get("fileSha1")
	fileName := r.Form.Get("fileName")
	// fileMeta := meta.GetFileMeta(fileSha1)
	fileMeta, _ := meta.GetFileMetaDB(fileSha1)
	oriFilePath := fileMeta.FilePath
	newFilePath := "./storage/tmp/" + fileName

	// Check whether the option type is legal
	if optionType != "0" {
		log.Fatal("Failed to access rights to rename files, err: \n" + "option type error occurs")

		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Update the data of a file
	if err := os.Rename(oriFilePath, newFilePath); err != nil {
		log.Fatal("Failed to rename file when updating, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the metadata of a file
	fileMeta.FileName = fileName
	fileMeta.FilePath = newFilePath
	fileMeta.UpdateAt = time.Now().Format("2006-01-02 15:04:05")
	// meta.SetFileMeta(fileSha1, fileMeta)
	_ = meta.SetFileMetaDB(fileSha1, fileMeta)

	// Return the http response
	data, err := json.Marshal(fileMeta)
	if err != nil {
		log.Fatal("Failed to convert fileMeta to JSON, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileDownloadHandler is to handle browser clients downloading files from the http server.
func FileDownloadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	fileSha1 := r.Form.Get("fileSha1")
	// fileMeta := meta.GetFileMeta(fileSha1)
	fileMeta, _ := meta.GetFileMetaDB(fileSha1)

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

// SingleFileQueryHandler is to handle querying files' metadata by fileSha1.
func SingleFileQueryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	fileSha1 := r.Form.Get("fileSha1")

	// Use file hash to query file's metadata
	// fileMeta := meta.GetFileMeta(fileSha1)
	fileMeta, _ := meta.GetFileMetaDB(fileSha1)
	data, err := json.Marshal(fileMeta)
	if err != nil {
		log.Fatal("Failed to convert fileMeta to JSON, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the http response
	w.Write(data)
}

// BatchFilesQueryHandler is to handle querying files' metadata by limitCount.
func BatchFilesQueryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	nickname := r.Form.Get("nickname")
	limitCount, _ := strconv.Atoi(r.Form.Get("limitCount"))

	// Use limit count to query files' metadata
	var fileMetas []meta.FileMeta
	userFileRecords, _ := db.GetUserFiles(nickname, limitCount)
	for _, userFileRecord := range userFileRecords {
		fileMeta, _ := meta.GetFileMetaDB(userFileRecord.FileSha1)
		fileMetas = append(fileMetas, fileMeta)
	}
	data, err := json.Marshal(fileMetas)
	if err != nil {
		log.Fatal("Failed to convert fileMetas to JSON, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the http response
	w.Write(data)
}

// FileDeleteHandler is to handle browser clients deleting files on the http server.
func FileDeleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	nickname := r.Form.Get("nickname")
	fileSha1 := r.Form.Get("fileSha1")
	// fileMeta := meta.GetFileMeta(fileSha1)
	fileMeta, _ := meta.GetFileMetaDB(fileSha1)

	// Delete the data of a file
	// if err := os.Remove(fileMeta.FilePath); err != nil {
	// 	log.Fatal("Failed to remove file when deleting, err: \n" + err.Error())

	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// Delete the metadata of a file
	// meta.DeleteFileMeta(fileSha1)
	_ = meta.DeleteFileMetaDB(fileSha1)
	_ = db.DeleteUserFile(nickname, fileSha1)

	// Return the http response
	data, err := json.Marshal(fileMeta)
	if err != nil {
		log.Fatal("Failed to convert fileMeta to JSON, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
