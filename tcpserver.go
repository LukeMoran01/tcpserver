package tcpserver

import (
	"fmt"
	"log"
	"net"
)

const (
	HOST     = "localhost"
	PORT     = "9001"
	TYPE     = "tcp"
	GREETING = "newplayer"
)

func logIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Server struct {
	connections []net.Conn
	listener    net.Listener
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
	for {
		var newConnection net.Conn
		newConnection, err = server.listener.Accept()
		logIfError(err)
		newId := server.addConnection(newConnection)
		go server.handleNewConnection(newId)
	}
}

func (server *Server) handleNewConnection(newId int) {
	fmt.Println("Connection made. Waiting for Greeting...")
	greeting := server.ReadAndTrimString(newId)
	newConnection := server.connections[newId]
	if greeting == GREETING {
		fmt.Println("Greeting accepted. Sending id:", newId)
		newConnection.Write([]byte("accepted"))
		newConnection.Write([]byte(fmt.Sprint(newId)))
	} else {
		fmt.Println("Greeting not accepted. Disconnecting..")
		newConnection.Write([]byte("disconnect"))
		newConnection.Close()
	}
}

func (server *Server) ReadAndTrimString(id int) string {
	return readAndTrimStringGeneric(server.connections[id])
}

func (server *Server) SendString(message string, id int) {
	sendStringGeneric(message, server.connections[id])
}

func (server *Server) addConnection(newConnection net.Conn) int {
	server.connections = append(server.connections, newConnection)
	connectionId := len(server.connections) - 1
	return connectionId
}

func (server *Server) CleanUpClosedConnections() {
	var err error
	fmt.Println("Cleaning..")
	for i := len(server.connections) - 1; i >= 0; i-- {
		err = ping(server.connections[i])
		if err != nil {
			fmt.Println("Cleaning up connection", i)
			server.connections = append(server.connections[:i], server.connections[i+1:]...)
		}
	}
}

func ping(connection net.Conn) error {
	response := make([]byte, 1)
	var err error
	_, err = connection.Read(response)
	return err
}
