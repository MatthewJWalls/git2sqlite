/*

This structure assists in unpacking git pack files.

usage:

    unpacker := Unpacker{somebytes}

    s := unpacker.String(4) // unpack 4 bytes as a string
    b := unpacker.Bytes(10) // unpack next 10 bytes
    n := unpacker.UInt(4)   // unpack next4 bytes as an unsigned int
    unpacker.Skip(4)        // skip forward 4 bytes

Unpacker assumes ints are Big Endian 32 bit numbers.

*/

package main

import (
	"fmt"
	"log"
	"bytes"
	"encoding/binary"
)

type Unpacker struct {
	raw []byte
	pos int
}

func (this *Unpacker) UInt(i int) int {
	var ret int32
	buf := bytes.NewBuffer(this.raw[this.pos:this.pos+i])
	err := binary.Read(buf, binary.BigEndian, &ret)

	if err != nil {
		log.Fatal(err)
	}

	this.pos += i
	return int(ret)
}

func (this *Unpacker) Skip(i int) {
	this.pos += i
}

func (this *Unpacker) Back(i int) {
	this.pos -= i
}

func (this *Unpacker) String(i int) string {
	val := fmt.Sprintf("%s", this.raw[this.pos:this.pos+i])
	this.pos = this.pos + i
	return val
}

func (this *Unpacker) Bytes(i int) []byte {
	val := this.raw[this.pos:this.pos+i]
	this.pos += i
	return val
}
