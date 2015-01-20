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

func TestIndexFileUnpack(t *testing.T) {

	log.Println("Unpacking index.")

	path := filepath.Join(
		"test-data", 
		"samplerepo", 
		"objects", 
		"pack", 
		"pack-f15a2879b1e01aece417a8c1d06f5079858633f4.idx",
	)

	rawbytes, readerr := ioutil.ReadFile(path)
			
	if readerr != nil {
		log.Fatal(readerr)
	}

	unpacker := Unpacker{raw:rawbytes}
	_ = unpacker

	log.Printf("index signature: %q", unpacker.Bytes(4))
	log.Printf("index version: %d", unpacker.UInt(4))

	for i := 0; i < 255; i += 1 {
		log.Printf("%d fanout: %d", i, unpacker.UInt(4))
	}

	size := unpacker.UInt(4)

	for i := 0; i < size; i += 1 {
		log.Printf("%d entry: %x", i, unpacker.String(20))
	}

	for i := 0; i < size; i += 1 {
		log.Printf("%d crc: %d", i, unpacker.UInt(4))
	}

	for i := 0; i < size; i += 1 {
		val := unpacker.UInt(4)
		log.Printf("%d offset: %d %b", i, val, val)
	}

}

func TestPackFileUnpack(t *testing.T) {

	log.Println("Unpacking pack file.")

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
