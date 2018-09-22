package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// Run the port scanner
func main() {
    var host, port string
    flag.StringVar(&host, "H", "", "host to scan")
    flag.StringVar(&port, "p", "", "port to scan")
    flag.Parse()

    if host == "" || port == "" {
        fmt.Println("Usage: portscan -H <host> -p port")
        os.Exit(1)
    }

    d := net.Dialer{Timeout: 20*time.Second}
    conn, errors := d.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
		fmt.Println(" ")
		fmt.Println("Running check for host", host, "on port", port)
		fmt.Println(" ")

    if errors != nil {
        if oerr, ok := errors.(*net.OpError); ok {
            switch oerr.Err.(type) {
                case *os.SyscallError:
                    fmt.Println("connect: connection refused to", host, "on port", port )
                default:
                    if oerr.Timeout() {
                        fmt.Println("connect: connection timed out to", host, "on port", port )
                    } else {
                        panic("Unknown connection error")
                    }
            }
        }
    } else {
    fmt.Println("connect: connection successful to", host, "on port", port )
    }

    if conn != nil {
        conn.Close()
    }
		fmt.Println(" ")
		
}
