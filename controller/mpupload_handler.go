package controller

import (
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"

	"../meta"
	db "../model"
	rPool "../model/redis"
	"../util"
)

// FileMPUploadInitHandler is to handle the initlization of every multipart upload try
func FileMPUploadInitHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	nickname := r.Form.Get("nickname")
	fileSha1 := r.Form.Get("fileSha1")
	fileSize, _ := strconv.Atoi(r.Form.Get("fileSize"))

	// Get a Redis pool connection
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// Prepare multipart upload info
	mpUploadMeta := meta.MPUploadMeta{
		FileSha1:   fileSha1,
		FileSize:   int64(fileSize),
		UploadID:   util.GenUploadID(nickname),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(fileSize) / (5 * 1024 * 1024))),
	}

	// Write multipart upload info to Redis
	rConn.Do("HSET", "MP_"+mpUploadMeta.UploadID, "file_sha1", mpUploadMeta.FileSha1)
	rConn.Do("HSET", "MP_"+mpUploadMeta.UploadID, "file_size", mpUploadMeta.FileSize)
	rConn.Do("HSET", "MP_"+mpUploadMeta.UploadID, "chunk_count", mpUploadMeta.ChunkCount)

	// Return the http response
	resp := util.RespMsg{
		Code: 0,
		Msg:  "SUCCESS",
		Data: mpUploadMeta,
	}
	w.Write(resp.JSONBytes())
}

// FileMPUploadPartHandler is to handle saving chunks
func FileMPUploadPartHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	uploadID := r.Form.Get("uploadID")
	chunkIndex := r.Form.Get("chunkIndex")

	// Get a Redis pool connection
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// Create a local file to store the uploaded file
	dirPath := "./storage/cache/" + uploadID + "/"
	filePath := dirPath + chunkIndex
	err := os.MkdirAll(dirPath, os.ModePerm)
	localFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Failed to create file when multipart uploading, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	// Write the file into physical disks
	buf := make([]byte, 1024*1024)
	for {
		cnt, err := r.Body.Read(buf)
		localFile.Write(buf[:cnt])
		if err != nil {
			break
		}
	}

	// Update multipart upload info in Redis
	rConn.Do("HSET", "MP_"+uploadID, "chunk_index_"+chunkIndex, 1)

	// Return the http response
	resp := util.RespMsg{
		Code: 0,
		Msg:  "SUCCESS",
		Data: nil,
	}
	w.Write(resp.JSONBytes())
}

// FileMPUploadFinishHandler is to handle merging chunks
func FileMPUploadFinishHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	uploadID := r.Form.Get("uploadID")
	nickname := r.Form.Get("nickname")
	fileSha1 := r.Form.Get("fileSha1")
	fileName := r.Form.Get("fileName")
	fileSize, _ := strconv.Atoi(r.Form.Get("fileSize"))

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

	// Get a Redis pool connection
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// Check whether all chunks have been uploaded
	res, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		log.Fatal("Failed to search redis when multipart finishing, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var totalCount int
	var savedCount int
	for i := 0; i < len(res); i += 2 {
		k := string(res[i].([]byte))
		v := string(res[i+1].([]byte))

		if k == "chunk_count" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chunk_index_") && v == "1" {
			savedCount++
		}
	}

	if totalCount != savedCount {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "FAILED",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
		return
	}

	// Merge all chunks then clear them
	desDirPath := "./storage/tmp"
	srcDirPath := "./storage/cache"
	if err := util.MergeChunks(desDirPath, srcDirPath, fileName, totalCount); err != nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "FAILED",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
		return
	}
	if err := util.ClearChunks(srcDirPath, totalCount); err != nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "FAILED",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
		return
	}

	// Define the metadata of a file
	fileMeta := meta.FileMeta{
		FileSha1:    fileSha1,
		FileName:    fileName,
		FileSize:    int64(fileSize),
		FilePath:    "./storage/tmp/" + fileName,
		EnableTimes: int64(enableTimes),
		EnableDays:  int64(enableDays),
		CreateAt:    time.Now().Format("2006-01-02 15:04:05"),
		UpdateAt:    time.Now().Format("2006-01-02 15:04:05"),
	}

	// Save the metadata of a file
	// meta.CreateFileMeta(fileMeta)
	_ = meta.CreateFileMetaDB(fileMeta)
	_ = db.AppendUserFile(nickname, fileMeta.FileSha1)

	// Return the http response
	resp := util.RespMsg{
		Code: 0,
		Msg:  "SUCCESS",
		Data: nil,
	}
	w.Write(resp.JSONBytes())
}

// FileMPUploadCancelHandler is to handle cancelling upload
func FileMPUploadCancelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	uploadID := r.Form.Get("uploadID")

	// Get a Redis pool connection
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// Get saved chunk count and clear all saved chunks
	res, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		log.Fatal("Failed to search redis when multipart finishing, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var savedCount int
	for i := 0; i < len(res); i += 2 {
		k := string(res[i].([]byte))
		v := string(res[i+1].([]byte))

		if strings.HasPrefix(k, "chunk_index_") && v == "1" {
			savedCount++
		}
	}

	srcDirPath := "./storage/cache"
	if err := util.ClearChunks(srcDirPath, savedCount); err != nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "FAILED",
			Data: nil,
		}
		w.Write(resp.JSONBytes())
		return
	}

	// Remove multipart upload info in Redis
	rConn.Do("DEL", "MP_"+uploadID)

	// Return the http response
	resp := util.RespMsg{
		Code: 0,
		Msg:  "SUCCESS",
		Data: nil,
	}
	w.Write(resp.JSONBytes())
}

// FileMPUploadStatusHandler is to handle the process of every multipart upload try
func FileMPUploadStatusHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the http request
	r.ParseForm()
	uploadID := r.Form.Get("uploadID")

	// Get a Redis pool connection
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// Calculate the ratio of saved chunk count and total chunk count
	res, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		log.Fatal("Failed to search redis when multipart finishing, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var totalCount int
	var savedCount int
	for i := 0; i < len(res); i += 2 {
		k := string(res[i].([]byte))
		v := string(res[i+1].([]byte))

		if k == "chunk_count" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chunk_index_") && v == "1" {
			savedCount++
		}
	}

	ratio := float64(savedCount) / float64(totalCount)

	// Return the http response
	resp := util.RespMsg{
		Code: 0,
		Msg:  "SUCCESS",
		Data: struct {
			Ratio float64
		}{
			Ratio: ratio,
		},
	}
	w.Write(resp.JSONBytes())
}
