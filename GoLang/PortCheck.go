package main

import (
        "os"
        "fmt"
        "net"
        "os/exec"
        "strings"
        "github.com/fatih/color"
)

func Hosts(cidr string) ([]string, error) {
        var numips int = 0
        var startip int = 0
        bits := strings.Split(cidr, "/")
        if len(bits) < 2 {
             fmt.Println("Invalid CIDR Notation, ", cidr)
             os.Exit(1)
        }

        requestedbits := bits[1]
        if (requestedbits == "31") {
           fmt.Println("CIDR format not valid, cannot use 31 bits")
           os.Exit(2)
        }

        ip, ipnet, err := net.ParseCIDR(cidr)
        if err != nil {
                return nil, err
        }

        var ips []string

        for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
                ips = append(ips, ip.String())
        }

        // remove network address and broadcast address

        numips = len(ips)
        if (requestedbits == "32") {
             numips = 1
        } else {
             numips--
             startip = 1
        }
        return ips[startip : numips], nil
}


func inc(ip net.IP) {
        for j := len(ip) - 1; j >= 0; j-- {
                ip[j]++
                if ip[j] > 0 {
                        break
                }
        }
}

func main() {
        mystring := " "
        if len(os.Args) < 2 {
            fmt.Print("No CIDR supplied \n")
            os.Exit(1)
        }
        CIDR := os.Args[1]
        fmt.Println(" ")
        mystring = "Pinging addresses for, "
        mystring += CIDR
        color.Yellow(mystring)
        hosts, _  := Hosts(CIDR)
        lastip := len(hosts)
        lastip--
        fmt.Println("First IP =", hosts[0])
        fmt.Println("Last IP  =", hosts[lastip])
        fmt.Println(" ")

        for i := 0; i < len(hosts); i++ {
            _, err := exec.Command("ping", "-c1", "-t3", "-w2", hosts[i]).Output()
            mystring = hosts[i]
            if err != nil {
                mystring += ", offline"
                color.Red(mystring)
            } else {
                mystring += ", online"
                color.Green(mystring)
            }
        }
}
