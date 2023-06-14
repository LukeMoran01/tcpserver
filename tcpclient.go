package tcpserver

import (
	"fmt"
	"net"
	"strconv"
)

type Client struct {
	connection net.Conn
	Player     Player
}

func (client *Client) CloseConnection() {
	client.connection.Close()
}

func (client *Client) ReadAndTrimString() string {
	return readAndTrimStringGeneric(client.connection)
}

func (client *Client) SendString(message string) {
	sendStringGeneric(message, client.connection)
}

func (client *Client) AskForId() int {
	client.SendString(GREETING)
	fmt.Println("Asking server for new id...")
	response := client.ReadAndTrimString()
	if response == "accepted" {
		fmt.Println("Greeting accepted. Receiving id...")
		response = client.ReadAndTrimString()
		id, err := strconv.Atoi(response)
		logIfError(err)
		return id
	} else {
		return -1
	}
}

func (client *Client) MakeConnection() {
	fmt.Println("Making connection to server...")
	tcpaddr, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	logIfError(err)
	client.connection, err = net.DialTCP(TYPE, nil, tcpaddr)
	logIfError(err)
}

func (client *Client) CreatePlayer() {
	id := client.AskForId()
	fmt.Println("Creating player with provided id:", id)
	client.Player.Id = id
}
