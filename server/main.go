package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
)

type message struct {
    senderID int
    contents string
}

type client struct {
    name string
    conn net.Conn
}

var (
    messages chan message = make(chan message)
    clients []client
)

func broadcast() {
    for {
        newestMessage := <-messages
        for clientID, client := range clients {
            fmt.Println("ClientID:", clientID)
            if newestMessage.senderID == clientID {
                continue
            }
            fmt.Fprintln(client.conn, newestMessage.contents)
        }
    }
}

func newActiveConnection(id int, userInfo *client) {
    from := bufio.NewScanner(userInfo.conn)
    defer userInfo.conn.Close()

    fmt.Fprintf(userInfo.conn, "Welcome to the server, %s!\n", userInfo.name)
    messages<- message{id, userInfo.name + " has joined the server."}

    var msgBuffer string
    for from.Scan() {
        msgBuffer = from.Text()
        fmt.Println("Got message:", msgBuffer)
        messages<- message{id, msgBuffer}
    }
    if err := from.Err(); err != nil {
        fmt.Println("Conneciton no.", id, "- Well, shit")
    }
}

func main() {
    // wow that was easy
    fmt.Println("Server starting...")
    server, err := net.Listen("tcp", ":6969")
    if err != nil {
        log.Fatalln("ERROR:", err)
    }

    // TODO: Give us a CLI to interact with the server
    // go func() {}()

    go broadcast()
    for clientNo := 0; true; clientNo++ {
        conn, err := server.Accept()
        if err != nil {
            fmt.Fprintln(os.Stderr, "Accept ERROR:", err)
            continue
        }
        clients = append(clients, client{"", conn})
        fmt.Println("Client created. ID:", clientNo)
        fmt.Println("Clients:", clients)

        // Each connection gets its own goroutine
        go newActiveConnection(clientNo, &clients[clientNo])
    }
}
