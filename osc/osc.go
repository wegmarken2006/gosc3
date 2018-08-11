package osc

import (
	"bytes"
	eb "encoding/binary"
	"fmt"
	"net"
	"time"
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

func EncodeF32(num float32) []byte {
	var buf []byte
	wbuf := bytes.NewBuffer(buf)
	eb.Write(wbuf, eb.BigEndian, float32(num))
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

func DecodeI8(buf []byte) int {
	var num uint8
	rbuf := bytes.NewReader(buf)
	eb.Read(rbuf, eb.BigEndian, &num)
	return int(num)
}

func DecodeI32(buf []byte) int {
	var num uint32
	rbuf := bytes.NewReader(buf)
	eb.Read(rbuf, eb.BigEndian, &num)
	return int(num)
}

func DecodeF32(buf []byte) float32 {
	var num float32
	rbuf := bytes.NewReader(buf)
	eb.Read(rbuf, eb.BigEndian, &num)
	return float32(num)
}

type portConfig struct {
	UdpIP   string
	UdpPort int
}

func align(n int) int {
	return 4 - n%4
}

func extend_(pad []byte, bts []byte) []byte {
	n := align(len(bts))
	outb := []byte{}
	outb = append(outb, bts...)
	for ind := 0; ind < n; ind++ {
		outb = append(outb, pad...)
	}
	return outb
}

func EncodeBlob(bts []byte) []byte {
	b1 := EncodeI32(len(bts))
	outb := []byte{}
	outb = append(outb, b1...)
	outb = append(outb, extend_([]byte{0}, bts)...)
	return outb
}

type IDatum interface{}

func EncodeDatum(dt IDatum) []byte {
	switch dt.(type) {
	case int:
		return EncodeI32((dt.(int)))
	case float32:
		return EncodeF32((dt.(float32)))
	case string:
		return EncodeStr((dt.(string)))
	case []byte:
		return EncodeBlob((dt.([]byte)))
	default:
		break
	}
	panic("enocdedatum")
}

func tag(dt IDatum) string {
	switch dt.(type) {
	case int:
		return "i"
	case float32:
		return "f"
	case string:
		return "s"
	case []byte:
		return "b"
	default:
		break
	}
	panic("tag")
}

func descriptor(id []IDatum) string {
	outs := ","
	for _, dt := range id {
		outs = outs + tag(dt)
	}
	return outs
}

type Message struct {
	Name   string
	LDatum []IDatum
}

func EncodeMessage(message Message) []byte {
	es := EncodeDatum(message.Name)
	ds1 := EncodeDatum(descriptor(message.LDatum))
	ds2 := []byte{}
	for _, elem := range message.LDatum {
		ds2 = append(ds2, EncodeDatum(elem)...)
	}
	es = append(es, ds1...)
	es = append(es, ds2...)
	return es
}
func SendMessage(message Message) {
	bmsg := EncodeMessage(message)
	OscSend(bmsg)
}

func ScStart() {
	OscSetPort()
	msg1 := Message{Name: "/notify", LDatum: []IDatum{1}}
	SendMessage(msg1)
	msg1 = Message{Name: "/g_new", LDatum: []IDatum{1, 1, 0}}
	SendMessage(msg1)
}

var pcfg portConfig

func OscSetPort() portConfig {
	pcfg.UdpIP = "127.0.0.1"
	pcfg.UdpPort = 57110
	return pcfg
}

func OscSend(message []byte) {
	m, _ := time.ParseDuration("2s")
	conn, err := net.DialTimeout("udp", pcfg.UdpIP, m)
	//defer conn.Close()

	if err != nil {
		panic("Connection error")
	}
	conn.Write(message)
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	fmt.Println(buff[:n])
}
