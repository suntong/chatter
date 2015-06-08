package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"chatter/sock"
	"golang.org/x/net/websocket"
	"encoding/json"
)

type Response struct {
	Status int32 `json:"status"`
	Command string `json:"command"`
	DocumentId string `json:"documentId"`
	Author string `json:"author"`
	Data string `json:"data"`
}

var pp = sock.Init()
var index = 0
var genericError, _ = json.Marshal(Response {
	400,
	"", "", "", "",
})

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
	http.Handle("/p2pAdd", websocket.Handler(p2pAddSocket))
	http.Handle("/p2pGet", websocket.Handler(p2pGetSocket))
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":22222", nil))

}

func p2pAddSocket(ws *websocket.Conn) {
	var doc string
	websocket.Message.Receive(ws, &doc)
	var jsonDecoded map[string]interface {}

	if err := json.Unmarshal([]byte(doc), &jsonDecoded); err != nil {
		fmt.Printf("Error decoding json, %s", err)
		websocket.Message.Send(ws, string(genericError))
		return
	}

	id := fmt.Sprintf("DOC%d", index)
	index++
	fmt.Printf("%s ->> %s\n", doc, id)
	pConfig := pp.AddNewPeer(ws, id)
	responseBody, err := json.Marshal(Response {
		200,
		"Join",
		id,
		"", "",
	})
	if err != nil {
		fmt.Printf("Error = %s", err)
	}

	fmt.Printf("Responding back %+v", string(responseBody))
	websocket.Message.Send(ws, string(responseBody))
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
