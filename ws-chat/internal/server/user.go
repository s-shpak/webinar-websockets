package server

import (
	"io"
	"sync"
)

type User struct {
	mux  *sync.Mutex
	conn io.ReadWriteCloser
	name string
	done chan struct{}
}

func NewUser(name string, conn io.ReadWriteCloser) *User {
	return &User{
		mux:  &sync.Mutex{},
		conn: conn,
		name: name,
		done: make(chan struct{}),
	}
}
