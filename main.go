package main

import (
	"log"

	"net/http"

	"github.com/julienschmidt/httprouter"

	"./controller"
)

func main() {
	router := httprouter.New()
	// Handle static resource
	router.ServeFiles("/static/*filepath", http.Dir("./static"))

	// Handle index direction
	router.GET("/", controller.IndexHandler)

	// Handle file operations
	router.POST("/file/upload", controller.FileUploadHandler)
	router.GET("/file/query/filehash/:fileHash", controller.SingleFileQueryHandler)
	router.GET("/file/query/limitcount/:limitCount", controller.BatchFilesQueryHandler)
	router.GET("/file/download/:fileHash", controller.FileDownloadHandler)
	router.PUT("/file/update", controller.FileUpdateHandler)
	router.DELETE("/file/delete", controller.FileDeleteHandler)

	// Handle user operations
	router.GET("/user/signup/get", controller.UserSignUpGetHandler)
	router.POST("/user/signup/post", controller.UserSignUpPostHandler)
	router.GET("/user/signin/get", controller.UserSignInGetHandler)
	router.POST("/user/signin/post", controller.UserSignInPostHandler)
	router.POST("/user/info", controller.GetUserInfoHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Failed to listen and serve, err: \n" + err.Error())
	}
}
