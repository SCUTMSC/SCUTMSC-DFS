package controller

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	db "../model"
	"../util"
)

// HTTPIntercepter is to do some verification before entering the handler
func HTTPIntercepter(h httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(
		func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Parse the http request
			r.ParseForm()
			nickname := r.Form.Get("nickname")
			token := r.Form.Get("token")

			// Verify nickname and token
			if len(nickname) < 3 || !isTokenValid(nickname, token) {
				resp := util.RespMsg{
					Code: int(util.StatusInvalidToken),
					Msg:  "FAILED",
					Data: nil,
				}
				w.Write(resp.JSONBytes())
			}

			// Enter the specific handler
			h(w, r, ps)
		},
	)
}

func isTokenValid(nickname string, token string) bool {
	// Check format
	if len(token) != 40 {
		return false
	}

	// Check duration
	if ts := token[32:]; util.Hex2Dec(ts) < time.Now().Unix()-86400 {
		return false
	}

	// Check consistency
	if !db.CheckUserToken(nickname, token) {
		return false
	}

	return true
}
