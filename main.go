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
	router.GET("/file/query/filesha1", controller.SingleFileQueryHandler)
	router.GET("/file/query/limitcount", controller.BatchFilesQueryHandler)
	router.GET("/file/upload", controller.FileDownloadHandler)
	router.PUT("/file/update", controller.FileUpdateHandler)
	router.DELETE("/file/delete", controller.FileDeleteHandler)

	router.POST("/file/ftupload", controller.FileFastUploadHandler)
	router.POST("/file/mpupload/init", controller.FileMPUploadInitHandler)
	router.POST("/file/mpupload/part", controller.FileMPUploadPartHandler)
	router.POST("/file/mpupload/finish", controller.FileMPUploadFinishHandler)
	router.POST("/file/mpupload/cancel", controller.FileMPUploadCancelHandler)
	router.POST("/file/mpupload/status", controller.FileMPUploadStatusHandler)
	router.GET("/file/attachment", controller.FileDownloadAttachmentHandler)
	router.GET("/file/url", controller.FileDownloadURLHandler)
	router.GET("/file/brdownload", controller.FileBRDownloadHandler)

	// Handle user operations
	router.POST("/user/signup", controller.UserSignUpHandler)
	router.POST("/user/signin", controller.UserSignInHandler)
	router.POST("/user/info", controller.AuthIntercepter(controller.GetUserInfoHandler))

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Failed to listen and serve, err: \n" + err.Error())
	}
}
