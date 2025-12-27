package wazero_pipe

import (
	"bytes"
	"context"
	_ "embed"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed test\.wasm
var testwasm []byte

func TestModule(t *testing.T) {
	var (
		ctx = context.Background()
		out = &bytes.Buffer{}
	)
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	hostModule := New()
	hostModule.Register(ctx, r)

	compiled, err := r.CompileModule(ctx, testwasm)
	if err != nil {
		panic(err)
	}
	cfg := wazero.NewModuleConfig().WithStdout(out)
	mod1, err := r.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		t.Errorf(`%v`, err)
		return
	}

	mod2, err := r.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		t.Errorf(`%v`, err)
		return
	}

	var meta *meta
	ctx, meta, err = hostModule.InitContext(ctx, mod1)
	if err != nil {
		t.Fatalf(`%v`, err)
	}
	if v := readUint32(mod1, meta.ptrBufCap); v != 1<<10 {
		t.Fatalf("incorrect buffer cap: %#v %d", meta, v)
	}

	// create pipe map
	pipes := make(map[uint32]chan []byte)
	ctx = context.WithValue(ctx, hostModule.ctxKey, pipes)

	var n uint64 = 12345
	go func() {
		_, err := mod1.ExportedFunction(`testSendUint64`).Call(ctx, 12345)
		if err != nil {
			t.Fatalf("%v\n%s", err, out.String())
		}
	}()
	stack, err := mod2.ExportedFunction(`testRecvUint64`).Call(ctx)
	if err != nil {
		t.Fatalf("%v\n%s", err, out.String())
	}
	if stack[0] != n {
		t.Fatalf("expected %d, got %d", n, stack[0])
	}
	hostModule.Stop()
}
