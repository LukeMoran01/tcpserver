package tcpserver

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	HOST      = "localhost"
	PORT      = "9001"
	TYPE      = "tcp"
	GREETING  = "newplayer"
	READY     = "ready"
	SEEN      = "seen"
	GAMESTART = "start"
	GIVEHAND  = "hand"
)

func logIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Server struct {
	connections []net.Conn
	listener    net.Listener
	game        Game
}

func (server *Server) StartServer() {
	var err error
	fmt.Println("Starting tcp server...")
	server.listener, err = net.Listen(TYPE, HOST+":"+PORT)
	logIfError(err)
}

func (server *Server) ListenForConnections() {
	defer server.listener.Close()
	var err error
	fmt.Println("Listening for connections...")
	for len(server.connections) < 2 {
		var newConnection net.Conn
		newConnection, err = server.listener.Accept()
		logIfError(err)
		newId := server.addConnection(newConnection)
		go server.handleNewConnection(newId)
		time.Sleep(time.Second)
	}
	fmt.Println("Player count met. Starting game.")
}

func (server *Server) handleNewConnection(newId int) {
	fmt.Println("Connection made. Waiting for Greeting...")
	greeting := server.ReadAndTrimString(newId)
	newConnection := server.connections[newId]
	if greeting == GREETING {
		fmt.Println("Greeting accepted. Sending id:", newId)
		newConnection.Write([]byte("accepted"))
		newConnection.Write([]byte(fmt.Sprint(newId)))
		ready := server.ReadAndTrimString(newId)
		if ready == READY {
			newConnection.Write([]byte("seen"))
		} else {
			newConnection.Write([]byte("invalid ready"))
			newConnection.Close()
		}
	} else {
		fmt.Println("Greeting not accepted. Disconnecting..")
		newConnection.Write([]byte("disconnect"))
		newConnection.Close()
	}
}

func (server *Server) ReadAndTrimString(id int) string {
	return readAndTrimStringGeneric(server.connections[id])
}

func (server *Server) addConnection(newConnection net.Conn) int {
	server.connections = append(server.connections, newConnection)
	connectionId := len(server.connections) - 1
	return connectionId
}

func (server *Server) StartGame() {
	server.game = CreateGame()
	server.game.RollHoard()
	for i, connection := range server.connections {
		fmt.Println("Contacting player", i)
		sendStringGeneric(GAMESTART, connection)
		response := readAndTrimStringGeneric(connection)
		if response != GIVEHAND {
			connection.Close()
		} else {
			fmt.Println("Player", i, "asked for hand.")
			fmt.Println("Sending player", i, "hand")
			server.SendHand(connection)
		}
	}
}

func (server *Server) SendHand(connection net.Conn) {
	hand := server.game.RollPlayerHand()
	connection.Write([]byte(strconv.Itoa(hand[0])))
	connection.Write([]byte(strconv.Itoa(hand[1])))
}

// To each connection send the length of the hoard and then each of the elements
func (server *Server) SendInitialHoard() {
	hoard := server.game.Hoard
	for _, connection := range server.connections {
		for _, element := range hoard {
			connection.Write([]byte(strconv.Itoa(element)))
		}
	}
}
