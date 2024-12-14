package main

import (
	"fmt"
	"flag"

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




// clears terminal
func clear_term() {
	fmt.Print("\033[H\033[2J")
}

// Query
func Query(name string, qtype int) (ip []byte) {

	if qtype == 1 { // A (IPv4)
		if name == "example.com" {
			ip =[]byte{93, 184, 216, 34}
		} else {
			ip = []byte{127, 0, 0, 1}
		}
	} else if qtype == 28 { // AAAA (IPv6)
		if name == "example.com" {
			ip =[]byte{0x26, 0x06, 0x28, 0x00, 0x02, 0x20, 0x00, 0x01,0x02, 0x48, 0x18, 0x93, 0x25, 0xc8, 0x19, 0x46}
		} else {
			ip = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
		}
	}
	//fmt.Printf("%s = %d", name,ip)
	return 
}


func main() {

	// Set Arguements 
	DNS_port := flag.Int("DNS Port", 53, "an int")

    flag.Parse()

	fmt.Printf("Server: 127.0.0.1:%d \n",*DNS_port)

	fmt.Println()

	//test value
	ip := Query("example.com", 1)

	fmt.Printf("IP: %d", ip)

}