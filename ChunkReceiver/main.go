// send back to ChunkReceiver later
package main

import (
	"FrameTypes"
	"UDPBroadcast"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var f *os.File

// Receiver of ChunkBasedFileTransferProtocol v1
func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.StringVar(port, "p", "8080", "Port to listen on")

	flag.Parse()

	address := fmt.Sprintf(":%s", *port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("Starting Server...\nDate: {%v}\nAddress %s\n", time.DateTime, address)
	fmt.Printf("Listening on port 8080\n")
	fmt.Printf("Successfully Started Server")

	go UDPBroadcast.BroadcastMessage()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
		}

		go handleConn(conn)

	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	//buff := make([]byte, 1024)
	var f *os.File
	meta := FrameTypes.FileMeta{}

	for {
		f, _, _, _ = readFrame(conn, f, &meta)
	}

}

func readFrame(r io.Reader, f *os.File, meta *FrameTypes.FileMeta) (file *os.File, typ byte, data []byte, err error) {
	header := make([]byte, 5)
	if _, err = io.ReadFull(r, header); err != nil {
		return
	}

	typ = header[0]
	//Well we know the next 4 bytes will be the length of the rest of the payload
	lenOfPayload := binary.BigEndian.Uint32(header[1:5])
	fmt.Printf("Length of Payload: %d\n", lenOfPayload)
	switch typ {
	case FrameTypes.FrameMeta:
		//Handle Metadata
		payload := handleMetadata(r, header)
		if err = json.Unmarshal(payload, meta); err != nil {
			panic(err)
		}

		fmt.Printf("Type: %x\n", FrameTypes.FrameMeta)
		fmt.Printf("Data: %s\n", string(payload))
		fmt.Printf("JSON: %v\n", meta)
		if f == nil {
			fmt.Printf("Opening file to be %s \n", meta.Name)
			f, err = os.Create(meta.Name)

			if err != nil {
				panic("Error opening file")
			}
		}
		//Now that we got that JSON we should probably store it and also reenfoce some rules
		return f, typ, payload, nil

	case FrameTypes.FrameData:
		payload := handleMetadata(r, header)
		fmt.Println("Sent: ", len(payload)/1000, " out of max size ", meta.Size)
		//fmt.Printf("Bytes: %v\n", string(payload))

		f.Write(payload)
		return f, typ, payload, nil

	case FrameTypes.FrameEOF:
		break
	}

	return f, typ, nil, nil
}

func handleMetadata(r io.Reader, header []byte) []byte {
	n := binary.BigEndian.Uint32(header[1:]) //Since it's the metadata read the rest of the length

	payload := make([]byte, n)

	if n > 0 {
		if _, err := io.ReadFull(r, payload); err != nil {
			return nil
		}
	}

	return payload
}

func paraseMetadataHeader() {
	//Enforce our rules here
	//1) version MUST be 1
	//Need to extract into a struct probably

}
