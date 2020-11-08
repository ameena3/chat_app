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

var wg sync.WaitGroup

func main() {
	times := 1000
	wg.Add(times)
	dialAndRead(times)
	wg.Wait()
}

func dialAndRead(i int) {
	for v := 0; v < i; v++ {
		go func(v int) {
			defer wg.Done()
			log.Println("Connecting to the server")
			conn, err := net.Dial("tcp", "localhost:8080")
			if err != nil {
				log.Printf("There was error in connecting to the server. The error is : %s.\n", err.Error())
				return
			}
			fmt.Fprintf(conn, fmt.Sprintf("Sending test message %d\r\n", v))
			rd := bufio.NewReader(conn)
			conn.SetReadDeadline(time.Time{})
			for {
				log.Println("reading from the server")
				val, err := rd.ReadString('\r')
				if err != nil {
					if err == io.EOF {
						log.Println("Server closed the connection")
						return
					}
					log.Printf("There was an error reading from the connection the error is %s.", err.Error())
					return
				}
				log.Println(val)
			}
		}(v)
	}
}
