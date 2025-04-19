package runtime

import (
	"os"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/bytecodealliance/wasmtime-go/v31"
	callbackqueue "github.com/dadamu/contract-wasmvm/internal/contract/callback-queue"
	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces"
	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces/testutil"
)

func BenchmarkRuntimeAddOne(b *testing.B) {
	// Setup suite with *testing.B
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	engine := wasmtime.NewEngine()
	// Read wasm from file
	wasmFile, err := os.ReadFile("testdata/test.wasm")
	if err != nil {
		b.Fatalf("failed to read module: %v", err)
	}

	module, err := wasmtime.NewModule(engine, wasmFile)
	if err != nil {
		b.Fatalf("failed to create module: %v", err)
	}

	repository := testutil.NewMockIContractRepository(ctrl)
	runtime := NewRuntimeFromModule(
		engine,
		callbackqueue.NewCallbackQueue(),
		&[]interfaces.ResultEvent{},
		repository,
		module,

		[]byte("state"),
		"contractId",
		100_000,
	)

	// Setup the mock repository
	loadedValue := []byte{1, 0, 0, 0}
	savedValue := []byte{2, 0, 0, 0}

	repository.EXPECT().
		LoadEntity("contractId", "test").
		Return(loadedValue).
		AnyTimes()

	repository.EXPECT().
		SaveEntity("contractId", "test", savedValue).
		AnyTimes()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := runtime.Run(interfaces.NewContractMessage("contractId", "addOne", []byte{}, "sender"))
		if err != nil {
			b.Fatalf("failed to run addOne: %v", err)
		}
	}
}
