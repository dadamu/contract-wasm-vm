package runtime

import (
	"os"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/bytecodealliance/wasmtime-go/v31"
	callbackqueue "github.com/dadamu/contract-wasmvm/internal/contract/callback-queue"
	repositorytestutil "github.com/dadamu/contract-wasmvm/internal/contract/repository/testutil"
)

func BenchmarkRuntimeAddOne(b *testing.B) {
	// Setup suite with *testing.B
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	engine := wasmtime.NewEngine()
	// Read wasm from file
	wasmFile, err := os.ReadFile("testdata/release.wasm")
	if err != nil {
		b.Fatalf("failed to read module: %v", err)
	}

	module, err := wasmtime.NewModule(engine, wasmFile)
	if err != nil {
		b.Fatalf("failed to create module: %v", err)
	}

	repository := repositorytestutil.NewMockIRepository(ctrl)
	runtime := NewRuntimeFromModule(
		callbackqueue.NewCallbackQueue(),
		engine,
		"test",
		repository,
		module,
		100_000,
	)

	// Setup the mock repository
	loadedValue := []byte{1, 0, 0, 0}
	savedValue := []byte{2, 0, 0, 0}

	repository.EXPECT().
		LoadEntity("test", "test").
		Return(loadedValue, nil).
		AnyTimes()

	repository.EXPECT().
		SaveEntity("test", "test", savedValue).
		Return(nil).
		AnyTimes()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := runtime.Run("addOne", []byte{}, []byte{})
		if err != nil {
			b.Fatalf("failed to run addOne: %v", err)
		}
	}
}
