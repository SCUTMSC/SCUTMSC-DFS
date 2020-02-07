package controller

import (
	"log"

	"io/ioutil"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Index is to handle directing to upload.html
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := ioutil.ReadFile("./static/view/upload.html")
	if err != nil {
		log.Fatal("Failed to read upload.html, err: \n" + err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileUpload is to handle browser clients uploading files to the http server.
func FileUpload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ...
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
