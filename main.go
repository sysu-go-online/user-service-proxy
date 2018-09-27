package main

import (
	"bytes"
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

	var method, uri, host string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &uri)
	h := strings.Split(string(b), "\n")[1]
	host = h[6:]

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

	// Parse subdomain from url
	DOMAINNAME := os.Getenv("DOMAIN_NAME")
	if len(DOMAINNAME) == 0 {
		DOMAINNAME = "localhost"
	}

	hostname := strings.Split(host[:len(host)-1], ".")
	if len(hostname) < 2 {
		tmp := strings.Split(hostname[1], ":")
		if tmp[0] != DOMAINNAME {
			fmt.Println("Invalid url hostname")
		}
	}
	url := "http://" + CONSULADDRESS + CONSULPORT + "/v1/kv/upstreams/"
	addr, err := model.GetValueWithKey(hostname[0], url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// send data to server
	fmt.Println(addr)
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
