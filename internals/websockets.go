package internals

import (
	"golang.org/x/net/websocket"
	"sync"
)

type Server struct {
	mu          sync.Mutex
	waitingConn *connection
}

type connection struct {
	ws      *websocket.Conn
	message chan string
}


func NewServer() *Server {
	return &Server{}
}

func (s *Server) HandleConnection(ws *websocket.Conn) {
	defer ws.Close()
	conn := &connection{
		ws:      ws,
		message: make(chan string),
	}

	s.mu.Lock()
	if s.waitingConn == nil {
		s.waitingConn = conn
		s.mu.Unlock()
	} else {
		peer := s.waitingConn
		s.waitingConn = nil
		s.mu.Unlock()

		go pairConnections(conn, peer)
	}

	for {
		var msg string
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			break
		}
		conn.message <- msg
	}
}

func pairConnections(conn1, conn2 *connection) {
    if err := websocket.Message.Send(conn1.ws, "Player 1"); err != nil {
        close(conn2.message)
        return
    }
    if err := websocket.Message.Send(conn2.ws, "Player 2"); err != nil {
        close(conn2.message)
        return
    }
	var wg sync.WaitGroup
	wg.Add(2)

	// Forward messages from conn1 to conn2
	go func() {
		defer wg.Done()
		for msg := range conn1.message {
			if err := websocket.Message.Send(conn2.ws, msg); err != nil {
				close(conn2.message)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for msg := range conn2.message {
			if err := websocket.Message.Send(conn1.ws, msg); err != nil {
				close(conn1.message)
				return
			}
		}
	}()

	wg.Wait()
	conn1.ws.Close()
	conn2.ws.Close()
}
