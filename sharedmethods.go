package tcpserver

import "net"

func readAndTrimStringGeneric(connection net.Conn) string {
	buffer := make([]byte, 16)
	messageLength, err := connection.Read(buffer)
	logIfError(err)
	return string(buffer[:messageLength])
}

func sendStringGeneric(message string, connection net.Conn) {
	_, err := connection.Write([]byte(message))
	logIfError(err)
}
