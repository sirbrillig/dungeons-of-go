package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func takeAction(client *Client, userInput string) {
	fmt.Println("Input: ", userInput)
	userInputNormal := strings.ToLower(userInput)

	quitExp := regexp.MustCompile(`\b(quit|exit|bye)\b`)
	if quitExp.MatchString(userInputNormal) == true {
		fmt.Println("closing connection for ", client.Name)
		client.sendMessage("you have fled the dungeon!")
		client.Close()
		return
	}

	lookExp := regexp.MustCompile(`\b(l|look)\b`)
	if lookExp.MatchString(userInputNormal) == true {
		client.sendMessage(client.CurrentRoom.Description)
		return
	}

	client.sendMessage("Sorry, I didn't understand what you said.")

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
	client.sendMessage("greetings and welcome to the dungeon!")

	dungeon := NewDungeon()
	client.Dungeon = dungeon
	client.CurrentRoom = dungeon.Rooms[Point{0, 0}]
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
