package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type (
	message struct {
		senderID int
		contents string
	}

	client struct {
		name    string
		conn    net.Conn
		scanner *bufio.Scanner
	}

	intSet map[int]struct{}
)

var (
	messages     = make(chan message)
	clients      = make(map[int]client) // switched from slices for O(1) deletion
	freedIndices = make(intSet)
	inSet        = struct{}{} // This should be const, but I can't make this const for some reason.

	quitSignal = make(chan struct{})
)

func displayHelp() {
	fmt.Println("Commands:")
}

func runActiveConnection(id int, userInfo *client) {
	defer userInfo.conn.Close()
	defer delete(clients, id)

	fmt.Fprintf(userInfo.conn, "Welcome to the server, %s!\n", userInfo.name)
	messages <- message{id, userInfo.name + " has joined the server."}

	var msgBuffer string
	for userInfo.scanner.Scan() {
		msgBuffer = userInfo.scanner.Text()
		fmt.Println("Got message:", msgBuffer) // Keeping him here to test encryption
		messages <- message{id, userInfo.name + ": " + msgBuffer}
	}
	if err := userInfo.scanner.Err(); err != nil {
		fmt.Println("Connection no.", id, "- Well, shit")
	}

	messages <- message{id, userInfo.name + " has left the server."}
	freedIndices[id] = inSet
}

func main() {
	fmt.Println("Server starting...")
	server, err := net.Listen("tcp", ":6969")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	// Listen for quit signal
	go func() {
		<-quitSignal // Blocks until it gets a quit message
		fmt.Println("Shutting down...")
		os.Exit(0) // The whole application. Not just this goroutine.
	}()

	fmt.Println("Type \"help\" or \"?\" for a list of commands")

	// CLI for server maintenance
	go func() {
		cli := bufio.NewScanner(os.Stdin)
		var cmd string
		for cli.Scan() {
			cmd = cli.Text()

			switch cmd {
			case "q", "quit":
				quitSignal <- struct{}{} // don't actually send any bytes, we just want to say "suspend/resume"
			case "h", "help", "?":
				displayHelp()
			default:
				fmt.Println("Unrecognized command. Type \"?\" for a list of commands.")
			}
		}
		if err := cli.Err(); err != nil {
			log.Fatalln("Unexpected CLI error, if you're reading this I fucked up pretty bad:", err)
		}
	}()

	// Broadcast messages to all clients as they come in
	go func() {
		for {
			newestMessage := <-messages
			for clientID, client := range clients {
				// TODO: remove when replace with ncurses UI for client
				if newestMessage.senderID == clientID {
					continue
				}
				fmt.Fprintln(client.conn, newestMessage.contents)
			}
		}
	}()

	// Accept new connections and get them set up
	for clientNo := 0; true; clientNo++ {
		conn, err := server.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Accept ERROR:", err)
			continue
		}

		// If a user previously disconnected and we have an unused index,
		// use the first available one.
		// Maybe replace with UUIDs?
		newID := -1
		for index := range freedIndices {
			if _, ok := freedIndices[index]; ok {
				newID = index
				break
			}
		}
		if newID == -1 {
			newID = clientNo
		}

		// Ask for name on the client's end, and if we got one without issue let the user in
		connScanner := bufio.NewScanner(conn)
		if connScanner.Scan() {
			name := connScanner.Text()
			clients[newID] = client{name, conn, connScanner}
			fmt.Println("Client created. ID:", newID)
			fmt.Println("Clients:", clients)

			// Each connection gets its own goroutine
			newClient := clients[newID] // why do maps make me do this?
			go runActiveConnection(newID, &newClient)
		} else {
			fmt.Fprintln(os.Stderr, "Failed to read name...")
		}
	}
}
