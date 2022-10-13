package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/AiredaleDev/logal-chat/client/ui"
)

func main() {
	ui.Setup()

	conn, err := net.Dial("tcp", ":6969")
	defer conn.Close()
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	messageBox := bufio.NewScanner(os.Stdin)

	// Listen for responses
	go func() {
		fromServer := bufio.NewScanner(conn)
		for fromServer.Scan() {
			reply := fromServer.Text()
			fmt.Println(reply)
		}

		if err := fromServer.Err(); err != nil {
			log.Fatalln("Listen ERROR:", err)
		}

		// we should be here if we got EOF
		fmt.Println("Server shut down.")
	}()

	// Send messages
	fmt.Print("Name: ")
	for messageBox.Scan() {
		message := messageBox.Text()
		if message == ":q" {
			break
		}
		fmt.Fprintln(conn, message)
	}

	if err := messageBox.Err(); err != nil {
		log.Fatalln("ERROR:", err)
	} else {
		fmt.Println("Quitting...")
	}
}
