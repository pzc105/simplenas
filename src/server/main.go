package main

import (
	"flag"
	"pnas/db"
	"pnas/log"
	"pnas/service"
	"pnas/setting"
)

func main() {
	var configPath = flag.String("c", "", "path of config file")
	setting.Init(*configPath)
	setting.InitDir()
	log.Init()
	db.Init()

	var ser service.CoreService
	ser.Init()
	ser.Serve()
}
