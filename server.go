package main
import ("bufio"
	"fmt"
	"net"
	"strings")
func main() {
	fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", ":8081") // listen on all interfaces on port 8081
	conn, _ := ln.Accept() // accept connection on port
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n') //read line
		fmt.Print("Message Received:", string(message))
		newMessage := strings.ToUpper(message) //change to upper
		conn.Write([]byte(newMessage + "\n")) //send upper string back
	}
}