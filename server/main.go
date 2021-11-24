package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
)

func main() {
    // wow that was easy
    fmt.Println("Server starting...")
    server, err := net.Listen("tcp", ":6969")
    if err != nil {
        log.Fatalln("ERROR:", err)
    }

    // TODO: Give us a CLI to interact with the server
    // go func() {}()

    for {
        conn, err := server.Accept()
        if err != nil {
            fmt.Fprintln(os.Stderr, "Accept ERROR:", err)
            continue
        }

        go func(conn net.Conn) {
            from := bufio.NewScanner(conn)
            var msgBuffer string
            for from.Scan() {
                msgBuffer = from.Text()
                fmt.Println("Got message:", msgBuffer)
                fmt.Fprintln(conn, "Echo:", msgBuffer)
            }
        }(conn)
    }
}
