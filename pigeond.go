package main

import (
	"flag"
	"fmt"

	log "pigeond/log"
	"pigeond/socket"
	"pigeond/tasks"
)

const socketFile = "/var/run/pigeond.socket"

func main() {

	// Config log
	logFile := flag.String("l", "/var/log/pigeond.log", "Pigeond log file")
	debug := flag.Bool("d", false, "Enable debug")
	flag.Parse()
	if err := log.ConfLog(*logFile, *debug); err != nil {
		panic(err)
	}
	if *debug == true {
		fmt.Println("Start pigeond server with debug=true, do't use it in production environment.")
	}

	// Launch server, TaskProxy as callback funcation
	if err := socket.LaunchServer("unix", socketFile, "", tasks.TaskProxy); err != nil {
		log.Log.Error(err)
		panic(err)
	}
}
