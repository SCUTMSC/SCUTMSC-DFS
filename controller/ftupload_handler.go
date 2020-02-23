package controller

import (
	"net/http"

	db "../model"
	"../util"
	"github.com/julienschmidt/httprouter"
)

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
