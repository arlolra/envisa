package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

var ENVISALINK string = "192.168.1.182:4025"

func checksum(line []byte) {
	ind := len(line) - 2
	var sum int64
	for _, b := range line[:ind] {
		sum += int64(b)
	}
	sent := string(line[ind:])
	calculated := strings.ToUpper(strconv.FormatInt(sum&0xFF, 16))
	if sent != calculated {
		log.Fatalf("Checksum doesn't match; %s %s", sent, calculated)
	}
}

func main() {
	conn, err := net.Dial("tcp", ENVISALINK)
	if err != nil {
		log.Fatalf("Could not connect; %v", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix {
			log.Fatal("Line too long")
		} else if err == io.EOF {
			log.Fatalf("Envisalink closed connection")
		} else if err != nil {
			log.Fatalf("Err reading line; %v", err)
		}
		command := string(line[:3])
		dataLen, ok := commands[command]
		if !ok {
			log.Printf("Unhandled command: %s\n", command)
			continue
		}
		data := string(line[3 : 3+dataLen])
		log.Println(command)
		log.Println(data)
		checksum(line)
	}
}
