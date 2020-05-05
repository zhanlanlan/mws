package main

import (
	"mws/router"
)

func main() {
	server := router.Route()
	server.Run(":8080")
}
