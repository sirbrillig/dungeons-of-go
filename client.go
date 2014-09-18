package main

import (
	"bytes"
	"fmt"
	"net"
)

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

func (c *Client) sendMessage(message string) {
	c.OutputChan <- []byte(message + "\n")
}
