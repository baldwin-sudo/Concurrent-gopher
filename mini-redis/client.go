package main

import (
	"bufio"
	"net"
	"strconv"
)

type Client struct {
	ID     string
	Conn   net.Conn
	Output chan string
	isDone bool
}

func initClient(Conn net.Conn, id int) *Client {
	return &Client{ID: "client-" + strconv.Itoa(id),
		Conn:   Conn,
		Output: make(chan string),
		isDone: false}
}
func (c *Client) ClientWrite() {
	w := bufio.NewWriter(c.Conn)
	defer close(c.Output)

	for o := range c.Output {

		_, err := w.WriteString(o)
		if err != nil {
			return // connection closed or error
		}
		if err := w.Flush(); err != nil {
			return // connection closed or error
		}

	}

}
