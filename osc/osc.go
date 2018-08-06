package osc

import (
	"bytes"
	eb "encoding/binary"
	//"fmt"
)

func EncodeI16(num int) []byte {
	var buf []byte
	wbuf := bytes.NewBuffer(buf)
	eb.Write(wbuf, eb.BigEndian, uint16(num))
	return wbuf.Bytes()
}

func EncodeI8(num int) []byte {
	var buf []byte
	wbuf := bytes.NewBuffer(buf)
	eb.Write(wbuf, eb.BigEndian, uint8(num))
	return wbuf.Bytes()
}

func EncodeI32(num int) []byte {
	var buf []byte
	wbuf := bytes.NewBuffer(buf)
	eb.Write(wbuf, eb.BigEndian, uint32(num))
	return wbuf.Bytes()
}

func EncodeStr(str string) []byte {
	byteArray := []byte(str)
	return byteArray
}

func StrPstr(str string) []byte {
	out := []byte{byte(len(str))}
	out = append(out, EncodeStr(str)...)
	return out
}

func DecodeI16(buf []byte) int {
	var num uint16
	rbuf := bytes.NewReader(buf)
	eb.Read(rbuf, eb.BigEndian, &num)
	return int(num)
}
