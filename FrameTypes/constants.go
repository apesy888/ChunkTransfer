package FrameTypes

type FileMeta struct {
	Version   uint8  `json:"version"`          // protocol version (1)
	Name      string `json:"name"`             // filename
	Size      int64  `json:"size"`             // total file size in bytes
	MTimeUnix int64  `json:"mtime_unix"`       // last modified time (unix seconds)
	Mode      uint32 `json:"mode,omitempty"`   // file mode (e.g. 0644 -> 420)
	ChunkSize uint32 `json:"chunk_size"`       // sender chunk size hint
	SHA256    string `json:"sha256,omitempty"` // optional whole-file checksum (hex)
}

const (
	FrameMeta byte = 0x01 // filename, size, mime,etc
	FrameData byte = 0x02 //raw file data byres (chunks)
	FrameEOF  byte = 0x03 //EOF
)
