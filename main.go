package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/sysu-go-online/user-service-proxy/model"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	b := make([]byte, 1024)
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	b = b[:n]
	// parse to golang http request
	r := bytes.NewReader(b)
	reader := bufio.NewReader(r)
	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Println(err)
		return
	}

	// Get host
	host := req.Host
	host = strings.Split(host, ".")[0]

	// get address from consul
	CONSULADDRESS := os.Getenv("CONSUL_ADDRESS")
	if len(CONSULADDRESS) == 0 {
		CONSULADDRESS = "localhost"
	}
	CONSULPORT := os.Getenv("CONSUL_PORT")
	if len(CONSULPORT) == 0 {
		CONSULPORT = "8500"
	}
	if CONSULPORT[0] != ':' {
		CONSULPORT = ":" + CONSULPORT
	}

	url := "http://" + CONSULADDRESS + CONSULPORT + "/v1/kv/upstreams/"
	addr, err := model.GetValueWithKey(host, url)
	if err != nil {
		log.Println(err)
		return
	}

	// send data to server
	server, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	server.Write(b)

	go io.Copy(server, client)
	io.Copy(client, server)
}
