package main

import (
	"log"
	"net"
	"strconv"
	"strings"
)

// Client

type Client struct {
	socket net.Conn
}

func (client *Client) init(ip string, port int) error {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	socket, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	client.socket = socket
	return nil
}

func (client *Client) send(message string) {
	client.socket.Write([]byte(message))
	client.socket.Write([]byte(StopCharacter))
}

func (client *Client) receive() []byte {
	buff := make([]byte, 1024)
	n, _ := client.socket.Read(buff)

	message := buff[:n]
	return message
}

func (client *Client) destroy() {
	if client.socket != nil {
		client.socket.Close()
	}
	client.socket = nil
}

// ClientHandle

func RunClient(verbose bool, ip string, port int, message string) int {
	client := new(Client)
	err := client.init(ip, port)
	if err != nil {
		log.Fatalf("%s", err)
		return 1
	}
	defer client.destroy()

	client.send(message)
	messageFromServer := client.receive()
	log.Printf("Receiving: %s", messageFromServer)
	return 0
}
