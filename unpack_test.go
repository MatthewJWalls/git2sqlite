package main

import (
	"testing"
	"encoding/binary"
)

func makeInteger(in int) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(in))
	return bs
}

func TestEverything(t *testing.T) {

	stringbytes := []byte{'t', 'e', 's', 't'}
	numberbytes := makeInteger(1337)
	allbytes := append(stringbytes, numberbytes...)

	unpacker := Unpacker{raw:allbytes}

	if unpacker.String(4) != "test" {
		t.Error("Did not read string properly")
	}

	if unpacker.UInt(4) != 1337 {
		t.Errorf("Did not read integer properly")
	}

}
