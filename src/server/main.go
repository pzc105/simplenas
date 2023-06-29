package main

import (
	"pnas/db"
	"pnas/log"
	"pnas/service"
	"pnas/setting"
)

func main() {
	setting.Init(".")
	setting.InitDir()
	log.Init()
	db.Init()

	var ser service.CoreService
	ser.Init()
	ser.Serve()
}
