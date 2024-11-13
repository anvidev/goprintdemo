package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mountRoutes(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: allowAllMiddleware(mux),
	}

	log.Println("starting print server")
	log.Fatal(srv.ListenAndServe())
}

func mountRoutes(m *http.ServeMux) {
	m.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
		return
	})
	m.HandleFunc("POST /print", printHandler)
}

type printRequest struct {
	Network string `json:"network"`
	Addr    string `json:"addr"`
	Data    string `json:"data"`
}

func printHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, errors.New("content-type must be set to application/json").Error(), http.StatusBadRequest)
		return
	}

	var payload printRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := newTCPConn(payload.Network, payload.Addr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(payload.Data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("print successful on addr %s\n", payload.Addr)
	w.Write([]byte("print successful"))
	return
}

func allowAllMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func newTCPConn(network, addr string) (*net.TCPConn, error) {
	tcp, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP(network, nil, tcp)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
