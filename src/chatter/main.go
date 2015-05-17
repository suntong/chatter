package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/satori/uuid"
	"net/url"
	"chatter/sock"
	"golang.org/x/net/websocket"
)

var server = sock.NewServer()

func submitNewPage(ws *websocket.Conn) {
	body, err := ioutil.ReadAll(ws.Request().Body)
	fmt.Printf(string(body))
	if err != nil {
		fmt.Printf("Unable to read: %s", err)
		ws.Write([]byte("Unable to read"))
	} else {
		fmt.Print(string(body))
		id := uuid.NewV1().String()
		newUrl := fmt.Sprintf("/sha?id=%s", id)
		fmt.Printf("Go to %s", newUrl)
	}
}

func openNewWebSocket(ws *websocket.Conn) {
	uri, _ := url.Parse(ws.Request().RequestURI)
	id, present := uri.Query()["id"]
	if present == false || len(id) != 1 {
		ws.Write([]byte("Bad Request"))
	} else {
		c := sock.NewClient(ws, server, id[0])
		ws.Write([]byte(c.Document))
	}
}

func main() {
	http.Handle("/createPage", websocket.Handler(submitNewPage))
	http.Handle("/sha", websocket.Handler(openNewWebSocket))
	http.Handle("/", http.FileServer(http.Dir("webroot")))
	log.Fatal(http.ListenAndServe(":22222", nil))

}
