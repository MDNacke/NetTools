package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	"strings"
	"github.com/fatih/color"
)

// Run the port scanner
func main() {
    var host, port string
		var PrintString string
    flag.StringVar(&host, "H", "", "host to scan")
    flag.StringVar(&port, "p", "", "port to scan")
    flag.Parse()

    if host == "" || port == "" {
        fmt.Println("Usage: portscan -H <host> -p port")
        os.Exit(1)
    }

    fmt.Println(" ")
    now := time.Now()
    startnow := now.Format("Mon Jan _2 15:04:05 2006")

    PrintString = fmt.Sprintf("%v - Running PortCheck for host %v on port %v", startnow, host, port)
    color.Green(PrintString)
    fmt.Println(" ")
    d := net.Dialer{Timeout: 20*time.Second}
    conn, errors := d.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
    now = time.Now()
    endnow := now.Format("Mon Jan _2 15:04:05 2006")
    if errors != nil {
        if oerr, ok := errors.(*net.OpError); ok {
            switch oerr.Err.(type) {
                case *os.SyscallError:
                    PrintString = fmt.Sprintf("%v - connect: connection refused to %v on port %v", endnow, host, port)
                default:
                    if oerr.Timeout() {
                        PrintString = fmt.Sprintf("%v - connection timed out to %v on port %v", endnow, host, port)
                    } else {
                        if strings.Contains(errors.Error(), "no such host") {
                            PrintString = fmt.Sprintf("%v - connect: no such host %v", endnow, host)
                        } else {
                            panic("Unknown connection error")
                        }
                    }
            }
        }
        	color.Red(PrintString)
    } else {
        PrintString = fmt.Sprintf("%v - connect: connection successful to %v on port %v", endnow, host, port)
        color.Green(PrintString)
    }

    if conn != nil {
        conn.Close()
    }
    fmt.Println(" ")
    fmt.Println("Ble,ble,ble that's all folks ")
}
