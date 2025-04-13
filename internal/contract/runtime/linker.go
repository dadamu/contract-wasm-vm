package runtime

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/bytecodealliance/wasmtime-go/v31"
	callbackqueue "github.com/dadamu/contract-wasmvm/internal/contract/callback-queue"
)

func (e *Runtime) prepareLinker() *wasmtime.Linker {
	linker := wasmtime.NewLinker(e.engine)

	// Add db.save function
	if err := linker.DefineFunc(e.store, "runtime", "db.save", e.saveEntry()); err != nil {
		panic(err)
	}

	// Add db.load function
	if err := linker.DefineFunc(e.store, "runtime", "db.load", e.loadEntry()); err != nil {
		panic(err)
	}

	// Add contract.call function
	if err := linker.DefineFunc(e.store, "runtime", "contract.call", e.callEntry()); err != nil {
		panic(err)
	}

	// Add env.abort function
	if err := linker.DefineFunc(e.store, "env", "abort", e.abortEntry()); err != nil {
		panic(err)
	}

	return linker
}

// `db.save` function that will be called from the WASM code
func (e *Runtime) saveEntry() func(caller *wasmtime.Caller, idPtr int32, dataPtr int32) {
	return func(caller *wasmtime.Caller, idPtr int32, dataPtr int32) {
		// TODO: consume fuel

		// Read the id string
		id := string(readBytes(caller, idPtr))

		var dataBz json.RawMessage = readBytes(caller, dataPtr)

		err := e.repository.SaveEntity(id, e.contractID, dataBz)
		if err != nil {
			panic(err)
		}
	}
}

// `db.load` function that will be called from the WASM code
func (e *Runtime) loadEntry() func(caller *wasmtime.Caller, idPtr int32) int32 {
	return func(caller *wasmtime.Caller, idPtr int32) int32 {
		// TODO: consume fuel

		// Read the id string
		id := string(readBytes(caller, idPtr))

		loaded, err := e.repository.LoadEntity(e.contractID, id)
		if err != nil && err == fmt.Errorf("contract %s with id %s not found", e.contractID, id) {
			// Write empty data to the memory if the entity is not found
			return writeBytes(caller, []byte{})
		} else if err != nil {
			panic(err)
		}

		return writeBytes(caller, loaded)
	}
}

// `contract.call` function that will be called from the WASM code
func (e *Runtime) callEntry() func(caller *wasmtime.Caller, contractIDPtr int32, methodPtr int32, argsPtr int32) {
	return func(caller *wasmtime.Caller, contractIDPtr int32, methodPtr int32, argsPtr int32) {
		// TODO: consume fuel

		// Read the contract ID and method name strings
		contractID := string(readBytes(caller, contractIDPtr))
		method := string(readBytes(caller, methodPtr))

		// Read the arguments
		args := readBytes(caller, argsPtr)

		// Call the contract method with the arguments
		e.callbackQueue.Enqueue(
			callbackqueue.NewCallbackMessage(
				contractID,
				method,
				args,
			),
		)
	}
}

// `env.abort` function that will be called from the WASM code
func (e *Runtime) abortEntry() func(caller *wasmtime.Caller, arg1, arg2, arg3, arg4 int32) {
	return func(caller *wasmtime.Caller, msgPtr, filePtr, line, column int32) {
		msg := readUTF16EncodedString(readBytes(caller, msgPtr))
		file := readUTF16EncodedString(readBytes(caller, filePtr))

		panic(fmt.Errorf("WASM called abort msg: %s, file: %s, line: %d, column: %d", msg, file, line, column))
	}
}

func readBytes(caller *wasmtime.Caller, ptr int32) []byte {
	memory := caller.GetExport("memory").Memory()
	memoryData := memory.UnsafeData(caller)
	length := binary.LittleEndian.Uint32(memoryData[ptr-4 : ptr])
	data := make([]byte, length)
	copy(data, memoryData[ptr:ptr+int32(length)])
	return data
}

func writeBytes(caller *wasmtime.Caller, data []byte) int32 {
	memory := caller.GetExport("memory").Memory().UnsafeData(caller)

	// Write the result to the memory
	malloc := caller.GetExport("__new").Func()
	if malloc == nil {
		panic("failed to find __new function")
	}

	// Allocate memory for the ArrayBuffer
	resultPtr, err := malloc.Call(caller, int32(len(data)), 1)
	if err != nil {
		panic(err)
	}

	offset := resultPtr.(int32)
	copy(memory[offset:offset+int32(len(data))], data)

	return offset
}

func readUTF16EncodedString(bz []byte) string {
	var decoded []byte
	for i := 0; i < len(bz); i += 2 {
		decoded = append(decoded, bz[i])
	}
	return string(decoded)
}
