package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type chunkType int

const (
	HeaderType chunkType = iota
	DataType
	EndType
	OtherType
)

type chunk struct {
	typDesc string
	typ     chunkType
	size    int
	data    []byte
	crc     uint32
}

// A PNG file must start with the following 8 bytes
var pngSignature = []byte{137, 80, 78, 71, 13, 10, 26, 10}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "PNG - PNG file reader")
		fmt.Fprintf(os.Stdout, "usage: png [path]")
		os.Exit(0)
	}

	fileName := os.Args[1]
	fmt.Fprintf(os.Stdout, "Reading from file '%s'\n", fileName)

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error opening file: %s\n", err)
		os.Exit(1)
	}

	// Read PNG signature from file
	buff, err := readBytes(f, 8)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error reading file : %s\n", err)
		os.Exit(1)
	}

	if compareBytes(buff, pngSignature) {
		fmt.Fprintf(os.Stdout, "%s is a valid png file!\n", fileName)
	} else {
		fmt.Fprintf(os.Stdout, "%s is not a png file!\n", fileName)
	}

	// Read chunks
	var c *chunk
	for {
		c, err = readChunk(f)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error reading file : %s\n", err)
			os.Exit(1)
		}

		// Report size of chunk
		fmt.Fprintf(os.Stdout, "Type        : %d\n", c.typ)
		fmt.Fprintf(os.Stdout, "Type (desc) : %s\n", c.typDesc)
		fmt.Fprintf(os.Stdout, "Size        : %d\n", c.size)
		fmt.Fprintf(os.Stdout, "--------------------\n")

		if c.typ == EndType {
			break
		}
	}

	f.Close()
}

func readChunk(f *os.File) (*chunk, error) {
	var c chunk

	// Read size of chunk
	buff, err := readBytes(f, 4)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk size : %w", err)
	}
	c.size = int(binary.BigEndian.Uint32(buff))

	// Read size of chunk
	buff, err = readBytes(f, 4)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk type : %w", err)
	}
	c.typDesc = string(buff)
	switch string(buff) {
	case "IHDR":
		c.typ = HeaderType
	case "IDAT":
		c.typ = DataType
	case "IEND":
		c.typ = EndType
	default:
		c.typ = OtherType
	}

	// Read data
	buff, err = readBytes(f, c.size)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk data : %w", err)
	}
	c.data = buff

	// Read CRC
	buff, err = readBytes(f, 4)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk CRC : %w", err)
	}
	c.crc = binary.BigEndian.Uint32(buff)

	return &c, nil
}

func reverseBytes(buffert []byte) {
	for i, j := 0, len(buffert)-1; i < j; i, j = i+1, j-1 {
		buffert[i], buffert[j] = buffert[j], buffert[i]
	}
}

func readBytes(f io.Reader, size int) ([]byte, error) {
	var buffer []byte
	buffer = make([]byte, size)

	_, err := f.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func compareBytes(buf1, buf2 []byte) bool {
	if len(buf1) != len(buf2) {
		return false
	}
	for i := range buf1 {
		if buf1[i] != buf2[i] {
			return false
		}
	}
	return true
}
