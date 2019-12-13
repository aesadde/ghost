package main

import "github.com/aesadde/ghost/ghost"

func main() {
	config := new(interface{})
	admin := ghost.NewDefaultAdminServer(config)
	server := ghost.NewDefaultServer()

	ghost.RunServers(admin, server)

}
