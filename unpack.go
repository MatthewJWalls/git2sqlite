/*

This structure assists in unpacking binary files.

usage:

    unpacker := Unpacker{somebytes}

    s := unpacker.String(4) // unpack 4 bytes as a string
    unpacker.Skip(10)       // skip next 10 bytes
    n := unpacker.Int(4)    // unpack 4 bytes as an int

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
	err := binary.Read(buf, binary.LittleEndian, &ret)

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
