package UDPBroadcast

import (
	"fmt"
	"net"
	"time"
)

func BroadcastMessage() {

	conn, _ := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP("192.168.1.255"), // 255.255.255.255
		Port: 9999,
	})

	//addr, _ := net.ResolveUDPAddr("udp4", "239.255.255.250:9999")
	//conn, _ := net.DialUDP("udp4", nil, addr)
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
