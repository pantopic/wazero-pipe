package pipe

import (
	"encoding/binary"
	"unsafe"
)

var (
	id     uint32
	bufCap uint32 = 1 << 10 // 1KB
	bufLen uint32
	buf    = make([]byte, int(bufCap))
	meta   = make([]uint32, 4)
)

//export __pipe
func __pipe() (res uint32) {
	meta[0] = uint32(uintptr(unsafe.Pointer(&id)))
	meta[1] = uint32(uintptr(unsafe.Pointer(&bufCap)))
	meta[2] = uint32(uintptr(unsafe.Pointer(&bufLen)))
	meta[3] = uint32(uintptr(unsafe.Pointer(&buf[0])))
	return uint32(uintptr(unsafe.Pointer(&meta[0])))
}

func setData(data any) (err error) {
	n, err := binary.Encode(buf[:bufCap], binary.LittleEndian, data)
	if err != nil {
		return
	}
	bufLen = uint32(n)
	return
}

func getData(data any) (err error) {
	_, err = binary.Decode(buf[:bufLen], binary.LittleEndian, data)
	return
}

//go:wasm-module pantopic/wazero-pipe
//export __host_pipe_recv
func recv()

//go:wasm-module pantopic/wazero-pipe
//export __host_pipe_send
func send()

// Fix for lint rule `unusedfunc`
var _ = __pipe
