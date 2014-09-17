package main

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// Client
// -----

type Client struct {
	Name       string
	InputChan  chan []byte
	OutputChan chan []byte
	Conn       net.Conn
	QuitChan   chan bool
}

func (c *Client) Read(buffer []byte) bool {
	bytesRead, error := c.Conn.Read(buffer)
	if error != nil {
		c.Close()
		fmt.Println(error)
		return false
	}
	fmt.Println("Read ", bytesRead, " bytes")
	return true
}

func (c *Client) Close() {
	c.QuitChan <- true
	c.Conn.Close()
}

func (c *Client) Equal(other *Client) bool {
	if bytes.Equal([]byte(c.Name), []byte(other.Name)) {
		if c.Conn == other.Conn {
			return true
		}
	}
	return false
}

// -----
// End Client

func takeAction(client *Client, userInput string) {
	fmt.Println("Input: ", userInput)
	userInputNormal := strings.ToLower(userInput)
	quitExp := regexp.MustCompile(`\b(quit|exit|bye)\b`)
	if quitExp.MatchString(userInputNormal) == true {
		fmt.Println("closing connection for ", client.Name)
		client.Close()
	}
}

// Wait for data to appear in a chan and call takeAction with the data
func handleUserInput(client *Client) {
	for {
		select {
		case userInput := <-client.InputChan:
			go takeAction(client, string(userInput))
			break
		}
	}
}

func handleUserOutput(userConn net.Conn, outputChan chan []byte) {
	for {
		select {
		case data := <-outputChan:
			go writeToConn(userConn, data)
			break
		}
	}
}

func writeToConn(userConn net.Conn, data []byte) {
	_, err := userConn.Write([]byte(data))
	if err != nil {
		fmt.Println("error while writing")
		return
	}
}

// Read on a Conn for data, then add that data to a chan
func clientReader(client *Client) {
	input := make([]byte, 2048)
	for client.Read(input) {
		client.InputChan <- input
	}
}

// Send a string to the outputChan
func sendOutput(data string, outputChan chan []byte) {
	dataBytes := []byte(data + "\n")
	outputChan <- dataBytes
}

func acceptAndMakeNewConnection(listener net.Listener) {
	inputChan := make(chan []byte)
	outputChan := make(chan []byte)

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("error making connection:", err)
		return
	}

	fmt.Println("connection made")

	client := &Client{Name: "foobar", InputChan: inputChan, OutputChan: outputChan, Conn: conn}
	go clientReader(client)
	go handleUserInput(client)
	go handleUserOutput(conn, outputChan)
	sendOutput("greetings and welcome to the dungeon!", outputChan)
}

// ----
// main
func main() {
	fmt.Println("starting up")

	listener, err := net.Listen("tcp", ":8889")
	if err != nil {
		fmt.Println("error starting server:", err)
		return
	}
	defer listener.Close()

	for {
		acceptAndMakeNewConnection(listener)
	}
}
