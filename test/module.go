package main

import (
	"github.com/pantopic/wazero-pipe/sdk-go"
)

var p *pipe.Pipe[uint64]

func main() {
	p = pipe.New[uint64]()
}

//export testSendUint64
func testSendUint64(n uint64) {
	err := p.Send(n)
	if err != nil {
		panic(err)
	}
}

//export testRecvUint64
func testRecvUint64() uint64 {
	res, err := p.Recv()
	if err != nil {
		panic(err)
	}
	return res
}

// Fix for lint rule `unusedfunc`
var _ = testSendUint64
var _ = testRecvUint64
