package main

import (
	"flag"
	"pnas/db"
	"pnas/log"
	"pnas/service"
	"pnas/setting"
	"pnas/video"
)

func main() {
	f := video.IsVideo("E:/s.mp4")
	print(f)
	var configPath = flag.String("c", "", "path of config file")
	setting.Init(*configPath)
	setting.InitDir()
	log.Init()
	db.Init()

	var ser service.CoreService
	ser.Init()
	ser.Serve()
}
