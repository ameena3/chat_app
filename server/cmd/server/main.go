package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var rc sync.Map

var lock sync.Mutex

func main() {
	// Listen
	log.Println("Strating the server")
	connection, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Printf("There was an error in starting the server %s", err.Error())
		return
	}
	log.Println("Listening for connections")
	// Accept
	for {
		conn, err := connection.Accept()
		if err != nil {
			log.Printf("There was an error in listening for the connection %s", err.Error())
			return
		}
		go handleConnections(conn)
	}
}

func handleConnections(conn net.Conn) {
	defer conn.Close()
	defer rc.Delete(conn.RemoteAddr())

	if _, ok := rc.Load(conn.RemoteAddr()); !ok {
		rc.Store(conn.RemoteAddr(), conn)
	}

	rd := bufio.NewReader(conn)
	conn.SetReadDeadline(time.Time{})
	log.Printf("handling connection for %v", conn.RemoteAddr())
	var val string
	var err error
	for {
		val, err = rd.ReadString(byte('\n'))
		if err != nil {
			if err == io.EOF {
				log.Println("Client closed the connection")
				return
			}
			log.Printf("There was an error in handling the connection n %s", err.Error())
			return
		}
		fmt.Println(val)
		body := `<!DOCTYPE HTML> THIS IS TEST </HTML> \n`
		fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
		fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
		fmt.Fprint(conn, "Content-Type: text/html\r\n")
		fmt.Fprint(conn, "\r\n")
		fmt.Fprintf(conn, "<HTML> THIS IS TEST </HTML> \n")
		rc.Range(func(key, c interface{}) bool {
			if c != conn && val != "" {
				wt := bufio.NewWriter(c.(net.Conn))
				log.Println("broadcasting")
				_, err = wt.WriteString(val)
				if err != nil {
					log.Printf("There was an error in writing to the connection %s", err.Error())
					return true
				}
				err = wt.Flush()
				if err != nil {
					log.Printf("There was an error in flushing to the connection %s", err.Error())
					return true
				}
				wt.Reset(c.(net.Conn))
			}
			return true
		})
	}
}
