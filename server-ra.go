package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Launching server...")
	//change the IP 192.168.1.133 to your IP
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("152.96.214.243"), Port: 8081})
	if err != nil {
		panic(err)
	}
	for {
		buffer := make([]byte, 1024)
		_, addr, _ := conn.ReadFromUDP(buffer)
		s := strings.Trim(string(buffer), "\x00")
		fmt.Printf("Message Received (len=%v): %v", len(s), s)
		newMessage := "hey this is your reply: " + strings.ToUpper(string(buffer)) //change to upper
		conn.WriteToUDP([]byte(newMessage+"\n"), addr)                             //send upper string back
	}
}
