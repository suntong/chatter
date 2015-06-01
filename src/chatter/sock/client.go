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

func NewClient(ws *websocket.Conn, s *Server, id string) Client {
	ch := make(chan(string))
	document, err := s.ReadDocumentContent(id)
	if err != nil {
		fmt.Printf("No content found for id=%s. Error=%s\n", id, err)
	}
	fmt.Printf("New client connected. Docuent=%s\n", document)
	return Client {
		*ws,
		s,
		id,
		document,
		ch,
	}
}

func (c *Client) Write(msg string) {
	select {
	case c.ch <- msg:
		fmt.Printf("Wrote to channel\n")
	default:
//		c.Server.Del(c)
		fmt.Printf("client %s is disconnected.\n", c.Id)

	}
}


//Request to write - to client
func (c Client) listenWrite() {
	for {
		select {
		case <- c.ch:
			var documentAsByte []byte
			len, err := c.Ws.Read(documentAsByte)
			if err != nil {
				fmt.Printf("Unable to read from socket: %s\n", err)
			} else {
				fmt.Printf("Read %d bytes\n", len)
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
				fmt.Printf("Unable to write to socket: %s\n", err)
			} else {
				fmt.Printf("Wrote %d bytes\n", len)
			}
		}
	}
}

func (c Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}
