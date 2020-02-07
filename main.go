package main

import (
	"log"

	"net/http"

	"github.com/julienschmidt/httprouter"

	"./controller"
)

func main() {
	router := httprouter.New()

	router.GET("/", controller.Index)
	router.POST("/dfs/upload", controller.FileUpload)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Failed to listen and serve, err: \n" + err.Error())
	}
}
