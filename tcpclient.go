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

func (client *Client) CreatePlayer() error {
	id := client.AskForId()
	fmt.Println("Creating player with provided id:", id)
	client.Player.Id = id
	client.Player.Hand = make([]int, 2)
	client.Player.HoardView = make([]int, 0, 5)
	client.Player.Money = 200
	client.SendString(READY)
	seen := client.ReadAndTrimString()
	if seen == SEEN {
		fmt.Println("Player set-up successful. Waiting game start..")
		return nil
	} else {
		return fmt.Errorf("problem setting up player")
	}
}

func (client *Client) PlayGame() {
	message := client.ReadAndTrimString()
	if message == GAMESTART {
		client.SendString(GIVEHAND)
		fmt.Println("Asking for hand")
		firstRoll := client.ReadAndTrimString()
		client.Player.Hand[0], _ = strconv.Atoi(firstRoll)
		secondRoll := client.ReadAndTrimString()
		client.Player.Hand[1], _ = strconv.Atoi(secondRoll)
	}
}

func (client *Client) GetInitialHoard() {
	for i := 0; i < 3; i++ {
		roll, _ := strconv.Atoi(client.ReadAndTrimString())
		client.Player.HoardView = append(client.Player.HoardView, roll)
	}
}
