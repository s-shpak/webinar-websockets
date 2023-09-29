package pinger

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gobwas/ws"
)

type Config struct {
	PingInterval       time.Duration
	LostPingsThreshold int
}

type Message struct {
	Frame ws.Frame
	Err   error
}

type Pinger struct {
	msgs chan Message
}

func NewPinger(cfg *Config) *Pinger {
	return &Pinger{
		msgs: make(chan Message, 1),
	}
}

func (p *Pinger) ProcessPing(ping ws.Frame) {
	payload := make([]byte, len(ping.Payload))
	if n := copy(payload, ping.Payload); n != len(ping.Payload) {
		p.msgs <- Message{
			Err: fmt.Errorf("failed to fully copy the received ping payload"),
		}
		return
	}
	pong := ws.Frame{
		Header: ws.Header{
			Fin:    true,
			OpCode: ws.OpPong,
		},
		Payload: payload,
	}
	p.msgs <- Message{
		Frame: pong,
	}
}

func (p *Pinger) SendPing() {
	payload := make([]byte, 10)
	if _, err := rand.Read(payload); err != nil {
		p.msgs <- Message{
			Err: fmt.Errorf("failed to generate a random payload: %w", err),
		}
		return
	}
	ping := ws.Frame{
		Header: ws.Header{
			Fin:    true,
			OpCode: ws.OpPing,
		},
		Payload: payload,
	}
	p.msgs <- Message{
		Frame: ping,
	}
}

func (p *Pinger) Messages() <-chan Message {
	return p.msgs
}
