package main

import (
	"strconv"

	"github.com/pantopic/wazero-pipe/sdk-go"
)

var (
	p0 = pipe.New[uint64]()
	p1 = pipe.New[uint32](pipe.WithID(1))
	p2 = pipe.New[[]byte](pipe.WithID(2))
)

func main() {}

//export testSendUint64
func testSendUint64(n uint64) {
	err := p0.Send(n)
	if err != nil {
		panic(err)
	}
}

//export testRecvUint64
func testRecvUint64() uint64 {
	res, err := p0.Recv()
	if err != nil {
		panic(err)
	}
	return res
}

//export testSendUint32
func testSendUint32(n uint64) {
	err := p1.Send(uint32(n))
	if err != nil {
		panic(err)
	}
}

//export testRecvUint32
func testRecvUint32() uint64 {
	res, err := p1.Recv()
	if err != nil {
		panic(err)
	}
	return uint64(res)
}

//export testSendBytes
func testSendBytes(n uint64) {
	err := p2.Send([]byte(strconv.Itoa(int(n))))
	if err != nil {
		panic(err)
	}
}

//export testRecvBytes
func testRecvBytes() uint64 {
	res, err := p2.Recv()
	if err != nil {
		panic(err)
	}
	i, err := strconv.Atoi(string(res))
	if err != nil {
		panic(err)
	}
	return uint64(i)
}

// Fix for lint rule `unusedfunc`
var _ = testSendUint64
var _ = testRecvUint64
var _ = testSendUint32
var _ = testRecvUint32
var _ = testSendBytes
var _ = testRecvBytes
