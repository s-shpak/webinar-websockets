package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	config "ws-chat/internal/config/client"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.GetConfig()
	dialer := ws.Dialer{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	conn, _, _, err := dialer.Dial(context.Background(), fmt.Sprintf("wss://:%s/connect", cfg.ServerPort))
	if err != nil {
		return fmt.Errorf("failed to dial the server: %w", err)
	}
	defer conn.Close()

	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	if err := chat(ctx, conn, cfg); err != nil {
		log.Println(err)
	}

	return nil
}

func chat(ctx context.Context, conn io.ReadWriteCloser, cfg config.Config) error {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go receiveMsg(ctx, wg, conn)
	go sendMsg(ctx, wg, conn, cfg.Msg)

	<-ctx.Done()

	if err := closeConnection(conn); err != nil {
		return fmt.Errorf("failed to close the ws connection: %w", err)
	}
	wg.Wait()

	log.Println("connection closed!")
	return nil
}

func receiveMsg(ctx context.Context, wg *sync.WaitGroup, conn io.ReadWriteCloser) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		h, r, err := wsutil.NextReader(conn, ws.StateClientSide)
		if err != nil {
			log.Printf("failed to read a frame: %v", err)
			continue
		}
		switch h.OpCode {
		case ws.OpText:
			msg, err := io.ReadAll(r)
			if err != nil {
				log.Printf("failed to read the received message: %v", err)
				continue
			}
			log.Printf("received message: %s", string(msg))
		case ws.OpClose:
			return
		default:
			log.Printf("unexpected frame with opcode %d", h.OpCode)
		}
	}
}

func sendMsg(ctx context.Context, wg *sync.WaitGroup, conn io.ReadWriteCloser, msg string) {
	defer wg.Done()

	for {
		w := wsutil.NewWriter(conn, ws.StateClientSide, ws.OpText)
		w.Write([]byte(msg))
		if err := w.Flush(); err != nil {
			log.Printf("failed to flush the writer buffer: %v", err)
		}
		log.Printf("message sent")

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * time.Duration(rand.Intn(4)+1)):
		}
	}
}

func closeConnection(conn io.ReadWriteCloser) error {
	if err := wsutil.WriteClientMessage(conn, ws.OpClose, nil); err != nil {
		return fmt.Errorf("failed to write a close frame: %w", err)
	}

	return nil
}
