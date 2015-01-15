package main

import (
	"testing"
	"log"
	"io/ioutil"
	"path/filepath"
	"encoding/binary"
)

func makeInteger(in int) []byte {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(in))
	return bs
}

func TestPackFileUnpack(t *testing.T) {

	path := filepath.Join(
		"test-data", 
		"samplerepo", 
		"objects", 
		"pack", 
		"pack-f15a2879b1e01aece417a8c1d06f5079858633f4.pack",
	)

	rawbytes, readerr := ioutil.ReadFile(path)
			
	if readerr != nil {
		log.Fatal(readerr)
	}

	unpacker := Unpacker{raw:rawbytes}

	signature := unpacker.String(4)
	version := unpacker.UInt(4)
	entries := unpacker.UInt(4)
	
	if signature != "PACK" {
		t.Error("Did not unpack correct signature.")
	}

	if version != 2 {
		t.Error("Did not unpack correct version.")
	}

	if entries <= 0 {
		t.Error("Did not unpack correct entry count.")
	}

	first := unpacker.Bytes(1)[0]
	finished := first >> 8 == 0
	otype := first << 1 >> 5

	log.Printf("first byte: %b", first)
	log.Printf("finished?: %s", finished)
	log.Printf("type: %b", otype)

	for finished = false; finished == false; {
		next := unpacker.Bytes(1)[0]
		finished = next >> 8 == 0
		log.Printf("next bytes: %b", next)
		log.Printf("finished?: %s", finished)
	}

}
