

package main

import(
	"fmt"
	"log"
	"net"
	"sync"
	"time"
	)



var Globalbuffer [1024]byte
var dmutex *sync.Mutex = new(sync.Mutex)
var dmutexMap *sync.Mutex = new(sync.Mutex)



type Client struct {
    Addr        *net.UDPAddr 
    lastActive time.Time
}

type ClientPool struct {
    p map[string]*Client
    m sync.Mutex
}

func (c ClientPool) Add(addr   *net.UDPAddr ) {
    c.m.Lock()
	
	fmt.Println("client ", addr)
	
    c.p[addr.String()] = &Client{
        Addr:       addr,
        lastActive: time.Now(),
    }
    c.m.Unlock()
}

func (c ClientPool) Remove(addr net.Addr) {
    c.m.Lock()
  	delete(c.p, addr.String())
    c.m.Unlock()
}





func dlockmap() {
	dmutexMap.Lock()
}

func dunlockmap() {
	dmutexMap.Unlock()
}

func dlock() {
	dmutex.Lock()
}

func dunlock() {
	dmutex.Unlock()
}


func handleServerUDPConnection(conn *net.UDPConn) {

	// here is where you want to do stuff like read or write to client
	
		
	
	for{
		    dlock()
			n, addr, err :=conn.ReadFromUDP(Globalbuffer[:])
			dunlock()


			fmt.Println("UDP  to Server : ", addr)
			fmt.Println("Received from UDP UDP client :  ", string(Globalbuffer[:n]))

			if err != nil{
				log.Fatal(err)
			}

				time.Sleep(1000)
	}

	

}
/*
connect ADD TO MAP
disconnect REMOVE FROM MAP

*/


func (c ClientPool) SendToAll(serverconn *net.UDPConn) {
    c.m.Lock()
    defer c.m.Unlock()
	sz := len(Globalbuffer)
     

    for k, client := range c.p {
        
       	fmt.Println("client ",client.Addr,sz)		
    	dlock()
		 err :=serverconn.WriteToUDP(Globalbuffer[:sz], client.Addr )
		 dunlock()
		  
		  if err != nil{
				log.Fatal(err)
			}
		  fmt.Println("key", k)
			
				
        }
    }

 
func BroadCastUDPtoClients(conn *net.UDPConn, serverconn *net.UDPConn) {

	// here is where you want to do stuff like read or write to client

	
	var buf[1024]byte
	

			cp := &ClientPool{
        p: map[string]*Client{},
    }
	
		// Tell the server to accept connections forever
		// and push new connections into the newConnections channel.
      	
		go func() {
		for {
			    n, cliaddr, err := conn.ReadFromUDP(buf[:])
     			fmt.Println("UDP client : ", cliaddr)
				fmt.Println("Received from UDP client :  ", string(buf[:n]))
				
				if err != nil{
				log.Fatal(err)
			}

				if ("CONNECT\r\n" == string(buf[:n])){
				    cp.Add(cliaddr)				
				}


        		if ("DISCONNECT\r\n" == string(buf[:n])){
				    cp.Remove(cliaddr)				
				}



		}
	}()
	
    
	//loop forever 
       for{
	 
     	 cp.SendToAll(serverconn)
	 
	 //may be small sleep
	
     	}
    
	}


	func CheckError(err error) {
		if err != nil{
			fmt.Println("Error: " , err)
			
		}
	}


	func main() {
	
	    hostName := "localhost"
		pIn := "6000"
		pOut := "6001"
		service_In := hostName + ":" + pIn
		service_Out := hostName + ":" + pOut

		udpAddr_In, err_in := net.ResolveUDPAddr("udp4", service_In)

		CheckError(err_in)

		udpAddr_Out, err_out := net.ResolveUDPAddr("udp4", service_Out)

		CheckError(err_out)



		// setup listener for incoming UDP connection
		ServerConn, err_in := net.ListenUDP("udp", udpAddr_In)



		ClientConn, err_Out := net.ListenUDP("udp", udpAddr_Out)

		CheckError(err_in)
		defer ServerConn.Close()

		CheckError(err_Out)
		defer ClientConn.Close()





		fmt.Println("UDP server up and listening on port 6000")
		fmt.Println("UDP Clients  Waiting  and listening on port 6001")

		



		// wait for UDP client to connect
	go  handleServerUDPConnection(ServerConn)
		//supposed to be a thread 
	   BroadCastUDPtoClients(ClientConn, ServerConn)


	}