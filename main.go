package main

import (
	"net"
	"flag"
	"log"
)

var targetAddr string
var listenAddr string

func init() {
	flag.StringVar(&targetAddr, "target", "", "transfer all connection / data to target addr")
	flag.StringVar(&listenAddr, "listen", ":8580", "address on which the transfer server will be listening")
}

func transfer(c1 net.Conn) {
	c2, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Println("failed to connect to target: ", err)
		c1.Close()
		return
	}
	go func() {
		buffer := make([]byte, 2048)
		for {
			n, err := c2.Read(buffer)
			if err != nil {
				log.Println("failed to read from target: ", err)
				c1.Close()
				c2.Close()
				break
			}
			_, err = c1.Write(buffer[:n])
			if err != nil {
				log.Println("failed to write to source: ", err)
				c2.Close()
				c1.Close()
				break
			}
		}
	}()
	buffer := make([]byte, 2048)
	for {
		n, err := c1.Read(buffer)
		if err != nil {
			log.Println("failed to read from source: ", err)
			c2.Close()
			c1.Close()
			break
		}
		_, err = c2.Write(buffer[:n])
		if err != nil {
			log.Println("failed to write to target: ", err)
			c1.Close()
			c2.Close()
			break
		}
	}
}

func main() {
	flag.Parse()
	c2, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Fatal("failed to connect to target: ", err)
	}
	c2.Close()
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("failed to listen on: ", listenAddr);
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go transfer(conn)
	}
}
