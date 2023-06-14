package main

import (
	"fmt"
	"os"

	"github.com/lukemoran01/tcpserver"
)

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "server" {
			server := tcpserver.Server{}
			server.StartServer()
			server.ListenForConnections()
		}
	} else {
		client := tcpserver.Client{}
		client.MakeConnection()
		client.CreatePlayer()
		fmt.Println(client.Player)
	}
}
