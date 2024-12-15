package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// For Coloured Text
var Reset = "\033[0m" 
var Red = "\033[31m" 
var Green = "\033[32m" 
var Yellow = "\033[33m" 
var Blue = "\033[34m" 
var Magenta = "\033[35m" 
var Cyan = "\033[36m" 
var Gray = "\033[37m" 
var White = "\033[97m"

func dns(port string) {

	port = ":" + port

	server, err := net.ResolveUDPAddr("udp4", port)
	if err != nil {
			fmt.Println(err)
			return
	}
	
	dns_server, err := net.ListenUDP("udp4", server)

	if err != nil {
		fmt.Printf("Port:%s%s%s is already use\n",Red,port,Reset)
		os.Exit(1)
	}
	defer dns_server.Close()

	for {
		buffer := make([]byte, 512)
		_, address, err := dns_server.ReadFromUDP(buffer)
		if err != nil {
			if network_error, ok := err.(net.Error); ok && network_error.Timeout() {
				continue
			}
			panic(err)
		}

		var id, flags, qdcount, ancount, nscount, arcount uint16
		binary.Read(bytes.NewReader(buffer[:12]), binary.BigEndian, &id)
		binary.Read(bytes.NewReader(buffer[2:]), binary.BigEndian, &flags)
		binary.Read(bytes.NewReader(buffer[4:]), binary.BigEndian, &qdcount)
		binary.Read(bytes.NewReader(buffer[6:]), binary.BigEndian, &ancount)
		binary.Read(bytes.NewReader(buffer[8:]), binary.BigEndian, &nscount)
		binary.Read(bytes.NewReader(buffer[10:]), binary.BigEndian, &arcount)

		// Make sure we're receiving a normal query
		if (flags&0xFE8F) != 0 || qdcount == 0 {
			continue
		}

		// Make result flags
		response_flags := uint16(0x8000)
		// If recursion is desired, set it as available
		if flags&0x0100 != 0 {
			response_flags |= 0x0180
		}

		// Read first query section
		query_name, offset := read_name(buffer, 12)
		var query_type, query_class uint16
		binary.Read(bytes.NewReader(buffer[offset:]), binary.BigEndian, &query_type)
		binary.Read(bytes.NewReader(buffer[offset+2:]), binary.BigEndian, &query_class)

		if query_class != 1 || (query_type != 1 && query_type != 28) {
			continue
		}
		response_name := make_name(query_name)
		response_data := query(query_name, query_type)

		// Encode result
		response := new(bytes.Buffer)
		binary.Write(response, binary.BigEndian, id)
		binary.Write(response, binary.BigEndian, response_flags)
		binary.Write(response, binary.BigEndian, uint16(1))
		binary.Write(response, binary.BigEndian, uint16(1))
		binary.Write(response, binary.BigEndian, uint16(0))
		binary.Write(response, binary.BigEndian, uint16(0))
		response.Write(response_name)
		binary.Write(response, binary.BigEndian, query_type)
		binary.Write(response, binary.BigEndian, query_class)
		response.Write([]byte{0xc0, 0x0c})
		binary.Write(response, binary.BigEndian, query_type)
		binary.Write(response, binary.BigEndian, query_class)
		binary.Write(response, binary.BigEndian, uint32(0))
		binary.Write(response, binary.BigEndian, uint16(len(response_data)))
		response.Write(response_data)
		dns_server.WriteToUDP(response.Bytes(), address)
	}
}

func query(query_name string, query_type uint16) []byte {
	var ip []byte
	if query_type == 1 { // A (ipv4)
		if query_name == "example.com" {
			ip = []byte{93, 184, 216, 34}
		} else {
			ip = []byte{127, 0, 0, 1}
		}
	} else if query_type == 28 { // AAAA (ipv6)
		if query_name == "example.com" {
			ip = []byte{0x26, 0x06, 0x28, 0x00, 0x02, 0x20, 0x00, 0x01, 0x02, 0x48, 0x18, 0x93, 0x25, 0xc8, 0x19, 0x46}
		} else {
			ip = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
		}
	}
	fmt.Println(query_name, "=", ip)
	return ip
}

func make_name(query_name string) []byte {
	var response bytes.Buffer
	for _, x := range bytes.Split([]byte(query_name), []byte{'.'}) {
		response.WriteByte(byte(len(x)))
		response.Write(x)
	}
	response.WriteByte(0)
	return response.Bytes()
}

func read_name(data []byte, offset int) (string, int) {
	var response string
	for {
		length := data[offset]
		offset++
		if length == 0 {
			break
		}
		if response != "" {
			response += "."
		}
		response += string(data[offset : offset+int(length)])
		offset += int(length)
	}
	return response, offset
}


func pop() {

	pop_listener, err := net.Listen("tcp", ":110")
	if err != nil {
		fmt.Printf("Port:%s110%s is already use\n",Red,Reset)
		os.Exit(1)
	}
	defer pop_listener.Close()


	pop_listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second))

	for {
		connection, err := pop_listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			panic(err)
		}

		connection.Write([]byte("+OK\r\n"))

		buffer := make([]byte, 1024)
		n, err := connection.Read(buffer)
		if err != nil {
			break
		}
		msg := string(buffer[:n])
		fmt.Print(msg)

		switch {
		case msg[:4] == "USER":
			connection.Write([]byte("+OK\r\n"))
		case msg[:4] == "PASS":
			connection.Write([]byte("+OK\r\n"))
		case msg[:5] == "STAT":
			connection.Write([]byte("+OK 0 0\r\n"))
		case msg[:4] == "QUIT":
			connection.Write([]byte("+OK\r\n"))
		}
	}
}

func main() {
	// Channel for goroutine syncs
	finish := make(chan bool)

	// Set Arguements 
	dns_port := flag.String("dport", "53", "Current Port used for DNS server")
    flag.Parse()
	
	// Run Everything
	
	// DNS Server
	go func() {
		dns(*dns_port)
	} ()

	// POP Server
	go func() {
		pop()
	} ()

	// Very Dirty Hack
	time.Sleep(100 * time.Millisecond)

	// Provide Information
	fmt.Printf("DNS Server: 127.0.0.1:%s%s%s \n",Green,*dns_port,Reset)

	<-finish // Routines to finish
}