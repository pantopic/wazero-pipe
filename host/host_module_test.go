package wazero_pipe

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed test\.wasm
var testwasm []byte

func TestModule(t *testing.T) {
	var (
		ctx = context.Background()
	)
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	hostModule := New()
	hostModule.Register(ctx, r)

	compiled, err := r.CompileModule(ctx, testwasm)
	if err != nil {
		panic(err)
	}
	cfg := wazero.NewModuleConfig().WithStdout(os.Stdout).WithName(`recv`)
	mod1, err := r.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		t.Errorf(`%v`, err)
		return
	}
	cfg = wazero.NewModuleConfig().WithStdout(os.Stdout).WithName(`send`)
	mod2, err := r.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		t.Errorf(`%v`, err)
		return
	}
	ctx, err = hostModule.InitContext(ctx, mod1)
	if err != nil {
		t.Fatalf(`%v`, err)
	}
	meta := get[*meta](ctx, hostModule.ctxKeyMeta)
	if v := readUint32(mod1, meta.ptrBufCap); v != 1<<10 {
		t.Fatalf("incorrect buffer cap: %#v %d", meta, v)
	}

	ctx = hostModule.ContextCopy(ctx, ctx)

	t.Run(`uint64`, func(t *testing.T) {
		go func() {
			for n := range 10 {
				_, err := mod1.ExportedFunction(`testSendUint64`).Call(ctx, uint64(n))
				if err != nil {
					panic(err.Error())
				}
			}
		}()
		for n := range 10 {
			stack, err := mod2.ExportedFunction(`testRecvUint64`).Call(ctx)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if stack[0] != uint64(n) {
				t.Fatalf("expected %d, got %d", n, stack[0])
			}
		}
	})
	t.Run(`uint32`, func(t *testing.T) {
		go func() {
			for n := range 10 {
				_, err := mod1.ExportedFunction(`testSendUint32`).Call(ctx, uint64(n))
				if err != nil {
					panic(err.Error())
				}
			}
		}()
		for n := range 10 {
			stack, err := mod2.ExportedFunction(`testRecvUint32`).Call(ctx)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if stack[0] != uint64(n) {
				t.Fatalf("expected %d, got %d", n, stack[0])
			}
		}
	})
	t.Run(`bytes`, func(t *testing.T) {
		go func() {
			for n := range 10 {
				_, err := mod1.ExportedFunction(`testSendBytes`).Call(ctx, uint64(n))
				if err != nil {
					panic(err.Error())
				}
			}
		}()
		for n := range 10 {
			stack, err := mod2.ExportedFunction(`testRecvBytes`).Call(ctx)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if stack[0] != uint64(n) {
				t.Fatalf("expected %d, got %d", n, stack[0])
			}
		}
	})

	hostModule.Stop()
}
