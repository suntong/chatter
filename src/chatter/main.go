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

type CreatePage struct{}
type OpenPage struct{}
var server = sock.NewServer()

func (h CreatePage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Unable to read: %s", err)
		w.Write([]byte("Unable to read"))
	} else {
		fmt.Print(string(body))
		id := uuid.NewV1().String()
		newUrl := fmt.Sprintf("/sha?id=%s", id)
		fmt.Printf("Redirecting to %s", newUrl)
		http.Redirect(w, r, newUrl, 302)
	}
}

func (h OpenPage) openNewWebSocket(ws *websocket.Conn) {
	uri, _ := url.Parse(r.RequestURI)
	id, present := uri.Query()["id"]
	if present == false || len(id) != 1 {
		w.Write([]byte("Bad Request"))
	} else {
		sock.NewClient(server)
		w.Write([]byte(id[0]))
	}
}

func main() {
	var submitNewPage CreatePage
	var openPage OpenPage
	http.Handle("/createPage", submitNewPage)
	http.Handle("/sha", websocket.Handler(openPage))
	http.Handle("/", http.FileServer(http.Dir("webroot")))
	log.Fatal(http.ListenAndServe(":22222", nil))

}
