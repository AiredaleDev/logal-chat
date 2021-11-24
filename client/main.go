package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
)

func main() {
    conn, err := net.Dial("tcp", ":6969")
    defer conn.Close()
    if err != nil {
        log.Fatalln("ERROR:", err)
    }

    messageBox := bufio.NewScanner(os.Stdin)
    reply := []byte{}
    for messageBox.Scan() {
        message := messageBox.Text()
        if message == ":q" {
            break
        }
        fmt.Fprintln(conn, message)
        conn.Read(reply)
        fmt.Println(string(reply))
    }

    if err := messageBox.Err(); err != nil {
        log.Fatalln("ERROR:", err)
    } else {
        fmt.Println("Quitting...")
    }
}
