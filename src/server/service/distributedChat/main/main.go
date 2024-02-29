package main

import (
	"flag"
	"fmt"
	distributedchat "pnas/service/distributedChat"
)

func main() {
	var role int
	flag.IntVar(&role, "role", 0, "")
	myAddress := flag.String("listen", "", "")
	myId := flag.String("myid", "", "path of config file")
	targetAddress := flag.String("taddress", "", "path of config file")
	targetId := flag.String("tid", "", "path of config file")
	slotRange := flag.String("slot", "", "path of config file")
	flag.Parse()
	var sSlot int
	var eSlot int
	fmt.Sscanf(*slotRange, "%d->%d", &sSlot, &eSlot)
	var c distributedchat.Cluster
	c.Start(&distributedchat.StartParams{
		Role:          distributedchat.NodeRole(role),
		MyId:          *myId,
		ListenAddress: *myAddress,
		TargetAddress: *targetAddress,
		TargetId:      *targetId,
		StartSlot:     sSlot,
		EndSlot:       eSlot,
	})
	c.Wait()
}
