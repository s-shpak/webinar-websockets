package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/gobwas/ws"

	serverConfig "ws-chat/internal/config/server"
)

func NewServer(cfg serverConfig.Config) (*http.Server, error) {
	srv, err := initHTTPServer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize an HTTP server: %w", err)
	}
	return srv, nil
}

func initHTTPServer(cfg serverConfig.Config) (*http.Server, error) {
	tlsCfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		InsecureSkipVerify: true,
	}
	if len(cfg.Port) == 0 {
		return nil, fmt.Errorf("server port must be specified")
	}
	chat := NewChat()
	return &http.Server{
		Addr:      fmt.Sprintf(":%s", cfg.Port),
		TLSConfig: tlsCfg,
		Handler:   NewRouter(chat),
	}, nil
}

type Router struct {
	*http.ServeMux
	chat *Chat
}

func NewRouter(chat *Chat) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		chat:     chat,
	}
	r.HandleFunc("/connect", r.connect)
	return r
}

func (rt *Router) connect(w http.ResponseWriter, r *http.Request) {
	if err := rt.handleConn(w, r); err != nil {
		log.Println(err)
	}
}

func (rt *Router) handleConn(w http.ResponseWriter, r *http.Request) error {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return fmt.Errorf("failed to upgrade the HTTP connection: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("failed to properly close a connection: %w", err)
		}
	}()

	done, err := rt.chat.Connect(conn)
	if err != nil {
		return fmt.Errorf("failed to connect a user: %w", err)
	}
	<-done

	return nil
}
