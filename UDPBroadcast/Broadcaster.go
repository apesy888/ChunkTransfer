package UDPBroadcast

import (
	"fmt"
	"net"
	"time"
)

func BroadcastMessage() {
	conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4bcast, // 255.255.255.255
		Port: 9999,
	})
	defer conn.Close()
	for {
		fmt.Println("Broadcasting data")
		BroadCast(conn)
		time.Sleep(5 * time.Second)
	}
}

func CreateConnection() (addr *net.UDPAddr, conn *net.UDPConn) {
	addr, _ = net.ResolveUDPAddr("udp4", "255.255.255.255:8888")
	conn, _ = net.ListenUDP("udp4", addr)

	return addr, conn
}

func BroadCast(conn *net.UDPConn) {

	conn.Write([]byte("Hello, Broadcast!"))
}
