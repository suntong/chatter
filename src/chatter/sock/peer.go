package sock

import (
	"golang.org/x/net/websocket"
	"fmt"
	"encoding/json"
)

type Peer struct {
	peerList map[string][]*PeerConfig
	AddChan chan(*PeerConfig)
}

type PeerConfig struct {
	ws *websocket.Conn
	channel chan(*PeerConfig)
	id string
}

func Init() (*Peer) {
	pList := make(map[string][]*PeerConfig)
	aChan := make(chan(*PeerConfig))
	return &Peer {
		pList,
		aChan,
	}
}

func (p *Peer) GetDebugData() string {
	debug := fmt.Sprintf("Total document count : %d\n", len(p.peerList))
	for k, v := range(p.peerList) {
		debug = fmt.Sprintf("%s documentId: %s connectedClient %s\n", debug, k, len(v))
	}
	return debug
}

func (p *Peer) AddNewPeer(w *websocket.Conn, id string) (*PeerConfig) {
	listenChannel := make(chan(*PeerConfig))
	config := &PeerConfig {
		w,
		listenChannel,
		id,
	}
	_, present := p.peerList[id]
	if present {

	}
//	go config.Listen()
	p.peerList[id] = append(p.peerList[id], config)
	return config
}

func (p *Peer) Listen() {
	for {
		select {
		case newclientAdd := <- p.AddChan:
				root := p.peerList[newclientAdd.id][0]
				fmt.Print("Root fount %s\n" , root.id)

				doc, err := root.GetDocument()
				if err != nil {
					fmt.Printf("Unable to read from root %s\n", err)
				} else {
					websocket.Message.Send(newclientAdd.ws, doc)
					fmt.Printf("Document read from root %s\n", doc)
				}
		}
	}
}

func (pConfig *PeerConfig) Listen() {
	for {
		select {
		case newClientConfig := <- pConfig.channel:
			fmt.Printf("HERER")
			doc, err := pConfig.GetDocument()
			fmt.Printf("Found document from root node %s\n", doc)
			if err != nil {
				newClientConfig.ws.Close()
			} else {
				websocket.Message.Send(newClientConfig.ws, doc)
			}
		}
	}
}

func (pConfig *PeerConfig) GetDocument() (string, error) {
	var document string
	err := websocket.Message.Send(pConfig.ws, getDocumentPayload(pConfig.id))
	if err != nil {
		fmt.Printf("Unable to send GetDocument command to root %s\n", err)
		return "", nil
	}
	err = websocket.Message.Receive(pConfig.ws, &document)
	if err != nil {
		fmt.Printf("Unable to read from socket: %s\n", err)
		return "", fmt.Errorf("Unable to read from socket: %s\n", err)
	} else {
		fmt.Printf("Read %s\n", document)
	}
	return document, nil
}

func getDocumentPayload(id string) (string) {
	data := make(map[string] interface{})
	data["documentId"] = id
	data["command"] = "GetDocument"
	payload, _ := json.Marshal(data)
	fmt.Printf("get document request payload %s", payload)
	return string(payload)
}
