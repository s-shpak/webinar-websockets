package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Chat struct {
	mux   *sync.Mutex
	users map[string]*User
}

func NewChat() *Chat {
	return &Chat{
		mux:   &sync.Mutex{},
		users: make(map[string]*User),
	}
}

func (c *Chat) Connect(conn io.ReadWriteCloser) (<-chan struct{}, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	u := c.createUser(conn)
	if _, ok := c.users[u.name]; ok {
		return nil, fmt.Errorf("user %s is already connected", u.name)
	}
	c.users[u.name] = u

	go c.handleUser(u)

	return u.done, nil
}

func (c *Chat) createUser(conn io.ReadWriteCloser) *User {
	return NewUser(fmt.Sprintf("user-%d", len(c.users)+1), conn)
}

func (c *Chat) handleUser(u *User) {
	defer c.removeUser(u)

	for {
		h, r, err := wsutil.NextReader(u.conn, ws.StateServerSide)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Printf("failed to read a frame from the user %s: %v", u.name, err)
			continue
		}
		switch h.OpCode {
		case ws.OpClose:
			if err := wsutil.WriteServerMessage(u.conn, ws.OpClose, nil); err != nil {
				log.Printf("failed to write a close frame: %v", err)
			}
			return
		case ws.OpText:
			if err := c.broadcast(r); err != nil {
				log.Printf("failed to broadcast the received message: %v", err)
			}
		default:
			log.Printf("unexpected frame with opcode %d", h.OpCode)
		}
	}
}

func (c *Chat) broadcast(r io.Reader) error {
	payload, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read the message: %w", err)
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	for _, u := range c.users {
		if err := wsutil.WriteServerText(u.conn, payload); err != nil {
			return fmt.Errorf("failed to send a message to the user %s: %w", u.name, err)
		}
	}
	return nil
}

func (c *Chat) removeUser(u *User) {
	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.users, u.name)
	close(u.done)
}
