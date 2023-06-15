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
			server.StartGame()
			server.SendInitialHoard()
			select {}
		}
	} else {
		client := tcpserver.Client{}
		client.MakeConnection()
		err := client.CreatePlayer()
		if err != nil {
			fmt.Println("Issue creating player. Closing connection.")
			client.CloseConnection()
		} else {
			client.PlayGame()
			client.GetInitialHoard()
			client.Player.ShowHoard()
			client.Player.ShowHand()
		}

	}
}
