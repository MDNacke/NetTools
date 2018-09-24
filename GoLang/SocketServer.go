package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "time"
    "github.com/fatih/color"
    "os/signal"
    "syscall"
)

func cleanup() {
    Now := time.Now()
    TimeNow := Now.Format("Mon Jan _2 15:04:05 2006")

    fmt.Println(" ")
    PrintString := fmt.Sprintf("%v - Is that all you got???????", TimeNow)
    color.Red(PrintString)

}

func main() {
    // Listen for incoming connections.
    var CONN_PORT string
    CONN_HOST := "0.0.0.0"
    CONN_TYPE := "tcp"
    flag.StringVar(&CONN_PORT, "p", "", "port to open")
    flag.Parse()

    if CONN_PORT == "" {
        fmt.Println("Usage: SocketServer -p port")
        os.Exit(1)
    }

    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()

    Now := time.Now()
    TimeNow := Now.Format("Mon Jan _2 15:04:05 2006")
    Listening := CONN_HOST + ":" + CONN_PORT
    MyString := fmt.Sprintf("%v - Listening on %v", TimeNow, Listening)
    color.Green(MyString)

// Setup a signal handler
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        cleanup()
        os.Exit(0)
    }()

// Start receiving connections
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn, CONN_PORT)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn, conn_port string) {
  var PrintString, IPAddress string
  Now := time.Now()
  TimeNow := Now.Format("Mon Jan _2 15:04:05 2006")

  fmt.Println(" ")

  // Make a buffer to hold incoming data.
  buf := make([]byte, 1024)
  // Read the incoming connection into the buffer.
  _, err := conn.Read(buf)
  if err != nil {
    if err.Error() != "EOF" {
        fmt.Println("Error reading:", err.Error())
    }
  }

  if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
    Now = time.Now()
    TimeNow = Now.Format("Mon Jan _2 15:04:05 2006")
    IPAddress = addr.IP.String()
    PrintString = fmt.Sprintf("%v - Connection made from %v on port %v", TimeNow, IPAddress, conn_port)
    color.Green(PrintString)
  }

  // Send a response back to person contacting us.
  conn.Write([]byte("Message received."))
  // Close the connection when you're done with it.
  conn.Close()
}
