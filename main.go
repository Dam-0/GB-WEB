package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

var active_port = `none`

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

func check_for_serial_port() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")

	} else if len(ports) == 1 {
		fmt.Println("Only " + Green + "1" + Reset + " Serial port found")
		active_port = ports[0]

	} else if len(ports) > 1 {
		ports_detail, err := enumerator.GetDetailedPortsList()
			if err != nil {
				log.Fatal(err)
				}

		fmt.Printf("Ports Found: " + Green + "(%v)\n" + Reset, len(ports))
		for _, port := range ports {
			fmt.Printf("- %v \n Extra: %v", port, ports_detail.Name)
		}
		fmt.Printf("Select a Port\n> ")

		for {

			fmt.Scan(&active_port) 

			active_port = strings.ToUpper(active_port)
			if slices.Contains(ports, active_port) {
				break
			} else {
				clear_term()
				fmt.Printf("Ports Found: " + Green + "(%v)\n" + Reset, len(ports))
				for _, port := range ports {
					fmt.Printf("- %v\n", port)
				}

				fmt.Println(Red + "Select a valid port" + Reset)
				fmt.Print("> ")
			}
		}
	}
}

func main() {
	clear_term()
	check_for_serial_port()

	fmt.Printf("Current Active Port: %v", active_port)

	//mode := &serial.Mode{
	//	BaudRate: 115200,
	//	Parity: serial.EvenParity,
	//	DataBits: 8,
	//	StopBits: serial.OneStopBit,
	//}

	//port, err := serial.Open(active_port, mode)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(port)
}