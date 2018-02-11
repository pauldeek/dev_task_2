

package main
 
import (
    "fmt"
    "net"
    "time"
    "strconv"
)
 
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}
 
func main() {
    ServerAddr,err := net.ResolveUDPAddr("udp","127.0.0.1:6000")
    CheckError(err)
 
    LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
    CheckError(err)
 
    Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    CheckError(err)
 
    defer Conn.Close()
    var buffer [1024] byte
	//buf := make([]byte, 1024)
	
	for i := 0; i < 1024; i++ {
	
	buffer[i]="a"
	}
	
    for {
        
        _,err := Conn.Write(buffer)
        if err != nil {
            fmt.Println(msg, err)
        }
        time.Sleep(1000)
    }
}