package main

import (
	"fmt"
	"io"
	"log"
	"net"
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
	fmt.Println(string(b))

	var method, host string
	// Get method and uri
	header := strings.Split(string(b), "\r\n")
	method = strings.Split(header[0], " ")[0]
	// uri = strings.Split(header[0], " ")[1]

	host = header[1][6:]

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

	hostname := strings.Split(host, ".")
	url := "http://" + CONSULADDRESS + CONSULPORT + "/v1/kv/upstreams/"
	addr, err := model.GetValueWithKey(hostname[0], url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// send data to server
	server, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Println("HTTP/1.1 200 Connection established")
	} else {
		server.Write(b)
	}

	go io.Copy(server, client)
	io.Copy(client, server)
}
