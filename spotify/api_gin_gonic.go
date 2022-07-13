/* TODO LIST:
1. Create API for fontend to communicate with backend
*/

package main

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jkravitz/mytrace"
)

const (
	tools    = "tools/"
	ui       = "ui/"
	play     = "play/"
	search   = "search/"
	device   = "device/"
	seek     = "seek/"
	volume   = "volume/"
	imageURL = "image/"
	getInfo  = "info/"
)

var (
	SERVER_URL = ""
	router     *gin.Engine
)

func initGin() (err error) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	var url string

	router.Static("/ui/", "../server_root/ui/build")
	router.Static("/ui_old/", "./server_root/ui/")
	router.Static("/dashboard/", "../server_root/dashboard/build")
	router.Static("/dashboard_old/", "./server_root/dashboard/build")
	// ui/play/?mode=search&query=all+about+that+bass

	//play
	url = filepath.Join(tools, play)
	router.GET(url, uiPlay)

	//search
	url = filepath.Join(tools, search)
	router.GET(url, uiSearch)

	//device access/mutating
	url = filepath.Join(tools, device)
	router.GET(url, uiDevice)

	//seek
	url = filepath.Join(tools, seek)
	router.GET(url, uiSeek)

	url = filepath.Join(tools, volume)
	router.GET(url, uiVolume)

	url = filepath.Join(tools, imageURL)
	router.GET(url, uiGetImageRef)

	url = filepath.Join(tools, getInfo)
	router.GET(url, uiGetInfo)

	//TODO: Serve wasm file

	// err = router.Run(":8081")
	err = router.RunTLS(":8081", "./key/ca-cert.pem", "./key/ca-key.pem")
	mytrace.Errhandle_Exit(err, err.Error())
	return
}
