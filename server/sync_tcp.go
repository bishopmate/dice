package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/bishopmate/dice/config"
)

func respond(command string, connection net.Conn) error {
	if _, err := connection.Write([]byte(command)); err != nil {
		return err
	}
	return nil
}

// TODO: Max read in one shot is 512 bytes
// To allow input > 512 bytes, do repeated read until
// we get EOF or designated delimiter
func readCommand(connection net.Conn) (string, error) {
	var buf []byte = make([]byte, 2048)
	// println("command", buf)
	n, err := connection.Read(buf[:])
	// println("n", n)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func RunSyncTCPServer() {
	log.Println("starting a synchronous TCP server on", config.Host, config.Port)

	var concurrent_clients int = 0

	// listening to the configured host port
	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {

		connection, err := lsnr.Accept()

		if err != nil {
			panic(err)
		}

		concurrent_clients += 1
		log.Println("client connected with address:", connection.RemoteAddr(), "concurrent clients", concurrent_clients)

		for {

			command, err := readCommand(connection)
			if err != nil {
				connection.Close()
				concurrent_clients -= 1
				log.Println("client disconnected", connection.RemoteAddr(), "concurrent clients", concurrent_clients)
				if err == io.EOF {
					log.Println("error while reading command", err)
					break
				}
				log.Println("error while reading command", err)
			}
			// println("command", command)
			log.Println("command", command)
			if err = respond(command, connection); err != nil {
				log.Print("error while writing:", err)
			}

		}

	}
}
