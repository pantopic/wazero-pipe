package wazero_pipe

import (
	"context"
	"log"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Name is the name of this host module.
const Name = "pantopic/wazero-pipe"

var (
	DefaultCtxKeyMeta  = `wazero_pipe_meta`
	DefaultCtxKeyPipes = `wazero_pipe_map`
)

type meta struct {
	ptrID     uint32
	ptrBufCap uint32
	ptrBufLen uint32
	ptrBuf    uint32
}

type hostModule struct {
	sync.RWMutex

	module      api.Module
	ctxKeyMeta  string
	ctxKeyPipes string
}

type Option func(*hostModule)

func New(opts ...Option) *hostModule {
	p := &hostModule{
		ctxKeyMeta:  DefaultCtxKeyMeta,
		ctxKeyPipes: DefaultCtxKeyPipes,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *hostModule) Name() string {
	return Name
}
func (p *hostModule) Stop() {}

// Register instantiates the host module, making it available to all module instances in this runtime
func (p *hostModule) Register(ctx context.Context, r wazero.Runtime) (err error) {
	builder := r.NewHostModuleBuilder(Name)
	register := func(name string, fn func(ctx context.Context, m api.Module, stack []uint64)) {
		builder = builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(fn), nil, nil).Export(name)
	}
	for name, fn := range map[string]any{
		"__host_pipe_send": func(ctx context.Context, pipe chan []byte, data []byte) {
			pipe <- append([]byte{}, data...)
		},
		"__host_pipe_recv": func(ctx context.Context, pipe chan []byte) []byte {
			return <-pipe
		},
	} {
		switch fn := fn.(type) {
		case func(ctx context.Context, pipe chan []byte, data []byte):
			register(name, func(ctx context.Context, m api.Module, stack []uint64) {
				meta := get[*meta](ctx, p.ctxKeyMeta)
				pipe, ok := p.pipes(ctx)[id(m, meta)]
				if !ok {
					pipe = make(chan []byte)
					p.pipes(ctx)[id(m, meta)] = pipe
				}
				fn(ctx, pipe, getData(m, meta))
			})
		case func(ctx context.Context, pipe chan []byte) []byte:
			register(name, func(ctx context.Context, m api.Module, stack []uint64) {
				meta := get[*meta](ctx, p.ctxKeyMeta)
				pipe, ok := p.pipes(ctx)[id(m, meta)]
				if !ok {
					pipe = make(chan []byte)
					p.pipes(ctx)[id(m, meta)] = pipe
				}
				setData(m, meta, fn(ctx, pipe))
			})
		default:
			log.Panicf("Method signature implementation missing: %#v", fn)
		}
	}
	p.module, err = builder.Instantiate(ctx)
	return
}

// InitContext retrieves the meta page from the wasm module
func (p *hostModule) InitContext(ctx context.Context, m api.Module) (context.Context, error) {
	stack, err := m.ExportedFunction(`__pipe`).Call(ctx)
	if err != nil {
		return ctx, err
	}
	meta := &meta{}
	ptr := uint32(stack[0])
	for i, v := range []*uint32{
		&meta.ptrID,
		&meta.ptrBufCap,
		&meta.ptrBufLen,
		&meta.ptrBuf,
	} {
		*v = readUint32(m, ptr+uint32(4*i))
	}
	return context.WithValue(ctx, p.ctxKeyMeta, meta), nil
}

// ContextCopy populates dst context with the meta page from src context.
func (h *hostModule) ContextCopy(src, dst context.Context) context.Context {
	dst = context.WithValue(dst, h.ctxKeyMeta, get[*meta](src, h.ctxKeyMeta))
	dst = context.WithValue(dst, h.ctxKeyPipes, make(map[uint32]chan []byte))
	return dst
}

func (p *hostModule) pipes(ctx context.Context) map[uint32]chan []byte {
	return get[map[uint32]chan []byte](ctx, p.ctxKeyPipes)
}

func get[T any](ctx context.Context, key string) T {
	v := ctx.Value(key)
	if v == nil {
		log.Panicf("Context item missing %s", key)
	}
	return v.(T)
}

func id(m api.Module, meta *meta) uint32 {
	return readUint32(m, meta.ptrID)
}

func readUint32(m api.Module, ptr uint32) (val uint32) {
	val, ok := m.Memory().ReadUint32Le(ptr)
	if !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
	return
}

func getData(m api.Module, meta *meta) (b []byte) {
	return read(m, meta.ptrBuf, meta.ptrBufLen, meta.ptrBufCap)
}

func dataBuf(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrBuf, 0, meta.ptrBufCap)
}

func setData(m api.Module, meta *meta, b []byte) {
	copy(dataBuf(m, meta)[:len(b)], b)
	writeUint32(m, meta.ptrBufLen, uint32(len(b)))
}

func read(m api.Module, ptrData, ptrLen, ptrMax uint32) (buf []byte) {
	buf, ok := m.Memory().Read(ptrData, readUint32(m, ptrMax))
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", ptrData, ptrLen)
	}
	return buf[:readUint32(m, ptrLen)]
}

func readUint64(m api.Module, ptr uint32) (val uint64) {
	val, ok := m.Memory().ReadUint64Le(ptr)
	if !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
	return
}

func writeUint32(m api.Module, ptr uint32, val uint32) {
	if ok := m.Memory().WriteUint32Le(ptr, val); !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
}
