package main

import (
	"fmt"
	"net"
)

func takeAction(userInput string) {
	fmt.Println("Input: ", userInput)
}

// Wait for data to appear in a chan and call takeAction with the data
func handleUserInput(inputChan chan []byte) {
	for {
		select {
		case userInput := <-inputChan:
			go takeAction(string(userInput))
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
func readFromConn(userConn net.Conn, dataChan chan []byte) {
	for {
		var input = make([]byte, 512)
		_, err := userConn.Read(input)
		if err != nil {
			fmt.Println("error while reading")
			return
		}
		dataChan <- input
	}
}

// Send a string to the outputChan
func sendOutput(data string, outputChan chan []byte) {
	dataBytes := []byte(data + "\n")
	outputChan <- dataBytes
}

func newConnection(listener net.Listener) {
	inputChan := make(chan []byte)
	outputChan := make(chan []byte)

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("error making connection:", err)
		return
	}

	fmt.Println("connection made")
	go readFromConn(conn, inputChan)
	go handleUserInput(inputChan)
	go handleUserOutput(conn, outputChan)
	sendOutput("greetings and welcome to the dungeon!", outputChan)
}

func main() {
	fmt.Println("starting up")

	listener, err := net.Listen("tcp", ":8889")
	if err != nil {
		fmt.Println("error starting server:", err)
		return
	}

	for {
		newConnection(listener)
	}
}
