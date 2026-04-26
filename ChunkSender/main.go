// send back to ChunkSender later
package main

import (
	"FrameTypes"
	"UDPBroadcast"
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// Sender of ChunkBasedFileTransferProtocol v1
func main() {
	//8192  bytes to send
	filePath := flag.String("file", "", "Path to the video file to send")
	flag.StringVar(filePath, "f", "", "Path to the video file to send")

	port := flag.String("port", "8080", "Port to listen on")
	flag.StringVar(port, "p", "8080", "Port to listen on")

	flag.Parse()

	//address := ":8080"
	message := "This is the secret message"

	//Now we need to get the address from the broadcast
	addr := UDPBroadcast.Listener()
	arr := strings.Split(addr, ":")
	hostName := arr[0]
	address := hostName + ":8080"

	fmt.Printf("Sending Data to Listener on address: {%s}\n", address)
	fmt.Printf("Sending file %s\n", *filePath)
	fileName := filepath.Base(*filePath)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer conn.Close()

	fileSize, err := getSize(*filePath)
	if err != nil {
		//log
		fmt.Println("Err getting fileSize: ", err)
	}

	jsonMetadata := jsonMetadataWithFileName(fileName, fileSize)

	fmt.Printf("Sending Message: { %s }", message)
	writeFrame(conn, []byte(jsonMetadata))
	//After frame start the loop to send the bytes
	_, err = os.ReadFile(*filePath)
	if err != nil {
		fmt.Println(err)
	}

	tempDir := os.TempDir()
	fmt.Println(tempDir)

	path := filepath.Join(tempDir, *filePath)
	f, _ := os.Create(path)
	defer f.Close()
	writeFilebytes(conn, *filePath)
	//Get the size in bytes. Depending on loop when we hit the max bytes to send we know when to stop

	//After the bytes have been sent finally send the EOF at the end

}

func writeFilebytes(w io.Writer, path string) {
	f, _ := os.Open(path)

	defer f.Close()
	buffer := make([]byte, 8192)
	r := bufio.NewReader(f)

	for {

		n, err := r.Read(buffer)

		if n > 0 {
			payload := buffer[:n]
			writeBody(w, payload)
		}

		if err == io.EOF {
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}

	}

}

func writeFile(w io.Writer, data []byte) {

	fmt.Println(len(data))
	totalSize := len(data)
	chunkSize := 8192
	bytes := data[:]

	//8192 loop for the size
	for i := 0; i < totalSize; i += chunkSize {
		fmt.Println(i + chunkSize)
		if i+chunkSize <= totalSize {
			//w.Write(bytes[i : i+chunkSize])
			writeBody(w, bytes[i:i+chunkSize])
		} else {
			writeBody(w, bytes[i:totalSize])
			//w.Write(bytes[i:totalSize])
		}

	}

	//writeFrame(w, part)

	//writeFrame(w, part)
}

func encode(w io.Writer) {
	//We need a type message
	//[FrameType][Length][paylaod]
	//FrameType HardCode for now and send it over

}

func writeFrame(w io.Writer, payload []byte) error {
	var header [5]byte
	header[0] = FrameTypes.FrameMeta
	binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))

	if err := writeFull(w, header[:]); err != nil {
		panic(err)
	}

	if len(payload) > 0 {
		if _, err := w.Write(payload); err != nil {
			return err
		}
	}
	return nil

}

func writeBody(w io.Writer, payload []byte) error {
	var header [5]byte
	header[0] = FrameTypes.FrameData
	binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))

	if err := writeFull(w, header[:]); err != nil {
		panic(err)
	}

	if len(payload) > 0 {
		if _, err := w.Write(payload); err != nil {
			return err
		}
	}
	return nil

}

func writeFull(w io.Writer, payload []byte) error {
	for len(payload) > 0 {
		n, err := w.Write(payload)
		if err != nil {
			return err
		}
		payload = payload[n:]
	}
	return nil
}

func getSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		// Handle error (e.g., file does not exist)
		//should we panic here if we can't find the file and change later?
		return 0, err
	}
	// Size returns length in bytes as an int64
	fmt.Printf("File size: %d bytes\n", fi.Size())
	return fi.Size(), nil
}

/*
	binary.Write(w, binary.BigEndian, FrameTypes.FrameMeta)
	binary.Write(w, binary.BigEndian, FrameTypes.FrameData)
	binary.Write(w, binary.BigEndian, FrameTypes.FrameEOF)
*/

func jsonMetadataWithFileName(fileName string, fileSize int64) string {
	return fmt.Sprintf("{\n  \"version\": 1,\n  \"name\": \"%s\",\n  \"size\": %d,\n  \"mtime_unix\": 1735512345,\n  \"mode\": 420,\n  \"chunk_size\": 262144,\n  \"sha256\": \"optional-hex-string\"\n}", fileName, fileSize)

}

const (
	jsonMetaData = "{\n  \"version\": 1,\n  \"name\": \"video.mp4\",\n  \"size\": 123456789,\n  \"mtime_unix\": 1735512345,\n  \"mode\": 420,\n  \"chunk_size\": 262144,\n  \"sha256\": \"optional-hex-string\"\n}"
)
