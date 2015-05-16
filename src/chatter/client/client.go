package client

import "golang.org/x/net/websocket"

type Client struct {

}

func (c Client) NewClient() *websocket.Conn {
	return nil
}
