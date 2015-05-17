package sock

import (
	"golang.org/x/net/websocket"
	"fmt"
)

type Client struct {
	Ws websocket.Conn
	Server *Server
	Document string
	Id string
	ch chan(string)
}

func NewClient(ws websocket.Conn, s *Server, id string) Client {
	ch := make(chan(string))
	document, err := s.ReadDocumentContent(id)
	if err != nil {
		fmt.Printf("No content found for id=%s. Error=%s", id, err)
	}
	fmt.Printf("New client connected. Docuent=%s", document)
	return Client {
		ws,
		s,
		id,
		document,
		ch,
	}
}

func (c *Client) Write(msg string) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.id)
		c.server.Err(err)
	}
}


//Request to write - to client
func (c Client) listenWrite() {
	for {
		select {
		default:
			var documentAsByte []byte
			len, err := c.Ws.Read(documentAsByte)
			if err != nil {
				fmt.Printf("Unable to read from socket: %s", err)
			} else {
				fmt.Printf("Read %d bytes", len)
				c.Document = string(documentAsByte)
			}
		}
	}
}

//Request to read - from server
func (c Client) listenRead() {
	for {
		select {
		default:
			len, err := c.Ws.Write([]byte(c.Document))
			if err != nil {
				fmt.Printf("Unable to write to socket: %s", err)
			} else {
				fmt.Printf("Wrote %d bytes", len)
			}
		}
	}
}

func (c Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}
