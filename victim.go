package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Launching server...")
	//change the IP 192.168.1.133 to your IP
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("192.168.1.133"), Port: 8082})
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n') //read line
		fmt.Printf("Message Received (len=%v): %v", len(message), string(message))
		newMessage := strings.ToUpper(message) //change to upper
		conn.Write([]byte(newMessage + "\n"))  //send upper string back
	}
}
