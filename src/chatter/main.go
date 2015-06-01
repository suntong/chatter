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
var pp = sock.Init()
var index = 0

func submitNewPage(ws *websocket.Conn) {
	body, err := ioutil.ReadAll(ws.Request().Body)
	fmt.Printf("%+v", string(body))
	if err != nil {
		fmt.Printf("Unable to read: %s", err)
		ws.Write([]byte("Unable to read"))
	} else {
		fmt.Printf("Body = %s <end>\n", string(body))
		fmt.Print(string(body))
		id := uuid.NewV1().String()
		c := sock.NewClient(ws, server, id)
		c.Write("test response")
//		go c.Listen()
		_, err := server.AddNewClient(id, c)

		if err != nil {
			fmt.Printf("Can not add client to the server %s", err)
		} else {
			newUrl := fmt.Sprintf("/sha?id=%s", id)
			fmt.Printf("Go to %s\n", newUrl)
		}
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
func createPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "webroot/createPage.html")
}

func debug(w http.ResponseWriter, r *http.Request) {
	uri, _ := url.ParseRequestURI(r.RequestURI)
	id := uri.Query()["id"]
	if len(id) == 0 {
		w.Write([]byte("Provide id"))
		return
	}
	d := pp.GetDebugData()
	fmt.Printf("%s\n", d)
}

func main() {
	go pp.Listen()
	http.HandleFunc("/debug", debug)
//	http.Handle("/sha", websocket.Handler(openNewWebSocket))
	http.Handle("/p2pAdd", websocket.Handler(p2pAddSocket))
	http.Handle("/p2pGet", websocket.Handler(p2pGetSocket))
	http.Handle("/", http.FileServer(http.Dir("webroot")))

//	http.HandleFunc("/createPage", createPage)
	log.Fatal(http.ListenAndServe(":22222", nil))

}

func p2pAddSocket(ws *websocket.Conn) {
	var doc string
	websocket.Message.Receive(ws, &doc)
	id := fmt.Sprintf("DOC%d", index)
	index++
	fmt.Printf("%s ->> %s\n", doc, id)
	pConfig := pp.AddNewPeer(ws, id)
	pConfig.Listen()
}

func p2pGetSocket(ws *websocket.Conn) {
	var id string
	err := websocket.Message.Receive(ws, &id)
	if err != nil {
		fmt.Printf("Unable to read id from client %s", err)
		return
	}
	fmt.Printf("Reading document with idx %s", id)
	config := pp.AddNewPeer(ws, id)
	fmt.Printf("Read complete for client %+v", config)
	pp.AddChan <- config
	config.Listen()
}
