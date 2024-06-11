package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /print/tsc", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		conn := newTSCTCP()
		defer conn.Close()

		if _, err := conn.Write(data); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("Print successful with"))
	})

	mux.HandleFunc("POST /print/hp", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		data, err := io.ReadAll(file)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		conn := newHPTCP()
		defer conn.Close()

		fmt.Println(data)
		bytesWritten, err := conn.Write(data)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(fmt.Sprintf("Print successful with %d bytes written", bytesWritten)))
	})

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	log.Println("Server is running")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func newTSCTCP() *net.TCPConn {

	tcp, err := net.ResolveTCPAddr("tcp4", "192.168.0.200:9100")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp4", nil, tcp)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func newHPTCP() *net.TCPConn {

	// addr := "2140:febb:acff:fa0d::fe80"
	tcp, err := net.ResolveTCPAddr("tcp4", "192.168.0.15:9100")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp4", nil, tcp)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
