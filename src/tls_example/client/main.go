package main

import (
	"crypto/tls"
	"fmt"
	"log"
)

// Simple TLS client
// Built from example at: https://gist.github.com/denji/12b3a568f092ab951456

func main() {
	log.SetFlags(log.Lshortfile)
	conf := &tls.Config{
	// InsecureSkipVerify: true
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	fmt.Println(string(buf[:n]))
}
