package main

import (
	"log"
	"testing"
	"encoding/binary"
)

func makeInteger(in int) []byte {
	bs := make([]byte, 4)
    binary.LittleEndian.PutUint32(bs, uint32(in))
	return bs
}

func TestEverything(t *testing.T) {

	log.Printf("%q", makeInteger(1337))

	stringbytes := []byte{'t', 'e', 's', 't'}
	numberbytes := makeInteger(1337)
	allbytes := append(stringbytes, numberbytes...)

	unpacker := Unpacker{raw:allbytes}

	if unpacker.String(4) != "test" {
		t.Error("Did not read string properly")
	}

	log.Printf("Unpacking from position %d", unpacker.pos)
	val := unpacker.UInt(4)

	//if unpacker.Int(1) != 2 {
	if val != 1337 {
		unpacker.Back(4)
		log.Printf("We unpacked: %q", unpacker.Bytes(4))
		t.Errorf("Did not read integer properly, was %d", val)
	}

}
