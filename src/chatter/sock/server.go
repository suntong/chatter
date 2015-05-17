package sock

import "fmt"

type Server struct {
	documentChannel map[string][]Client
}

func NewServer() *Server {
	documentChannel := make(map[string][]Client)
	return &Server{
		documentChannel,
	}
}

func (s *Server) AddNewClient(id string, client Client) (bool, error) {
	clientChannel, present := s.getDocument(id)
	if !present {
		fmt.Printf("No such document found: %s", id)
		return false, fmt.Errorf("No such document found: %s", id)
	}
	append(clientChannel, client)
	return true, nil
}

func (s *Server) DocumentViewCount(id string) int {
	clientChannel, present := s.getDocument(id)
	if !present {
		return 0
	}
	return len(clientChannel)
}

func (s *Server) ReadDocumentContent(id string) (string, error) {
	clientChannel, present := s.getDocument(id)
	if !present {
		fmt.Printf("No such document found: %s", id)
		return nil, fmt.Errorf("No such document found: %s", id)
	} else {
		return clientChannel[0].Document, nil
	}
}

func (s *Server) WriteDocumentContent(c Client, document string) {
	len, err := c.Ws.Write([]byte(document))
	if err != nil {
		fmt.Printf("Unable to write to client %s", err)
		return
	} else {
		fmt.Printf("Wrote %d byte to clinet", len)
	}
}

func (s *Server) getDocument(id string) (clientChannel []Client, present bool) {
	clientChannel, present := s.documentChannel[id]
	return
}

