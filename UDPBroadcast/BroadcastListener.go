package UDPBroadcast

import (
	"fmt"
	"net"
)

func Listener() string {
	// Listener example (brief)
	addr := net.UDPAddr{
		Port: 9999,
		IP:   net.IPv4zero, // 0.0.0.0 = listen on all interfaces
	}

	conn, err := net.ListenUDP("udp4", &addr)
	//addr, _ := net.ResolveUDPAddr("udp4", "0.0.0.0:9999")
	//conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Listening for UDP broadcast on port 9999...")

	buffer := make([]byte, 1024)

	//for {
	n, remoteAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	fmt.Printf("Received %d bytes from %s: %s\n",
		n,
		remoteAddr.String(),
		string(buffer[:n]),
	)
	return remoteAddr.String()
	//}

}

func setUpListener() (addr *net.UDPAddr, conn *net.UDPConn) {
	addr, _ = net.ResolveUDPAddr("udp4", "0.0.0.0:8888")
	conn, _ = net.ListenUDP("udp4", addr)

	return addr, conn
}
