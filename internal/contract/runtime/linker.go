package runtime

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/bytecodealliance/wasmtime-go/v31"

	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces"
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

	// Add contract.create function
	if err := linker.DefineFunc(e.store, "runtime", "contract.create", e.createContractEntry()); err != nil {
		panic(err)
	}

	// Add event.emit function
	if err := linker.DefineFunc(e.store, "runtime", "event.emit", e.emitEventEntry()); err != nil {
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
		id := readUTF16EncodedString(readBytes(caller, idPtr))

		var dataBz = readBytes(caller, dataPtr)
		e.repository.SaveEntity(e.contractId, id, dataBz)
	}
}

// `db.load` function that will be called from the WASM code
func (e *Runtime) loadEntry() func(caller *wasmtime.Caller, idPtr int32) int32 {
	return func(caller *wasmtime.Caller, idPtr int32) int32 {
		// TODO: consume fuel

		// Read the id string
		id := readUTF16EncodedString(readBytes(caller, idPtr))

		loaded := e.repository.LoadEntity(e.contractId, id)
		return writeBytes(caller, loaded)
	}
}

// `contract.call` function that will be called from the WASM code
func (e *Runtime) callEntry() func(caller *wasmtime.Caller, contractIdPtr int32, methodPtr int32, argsPtr int32) {
	return func(caller *wasmtime.Caller, contractIdPtr int32, methodPtr int32, argsPtr int32) {
		// TODO: consume fuel

		// Read the contract Id, method name strings and args
		contractId := readUTF16EncodedString(readBytes(caller, contractIdPtr))
		method := readUTF16EncodedString(readBytes(caller, methodPtr))
		args := readBytes(caller, argsPtr)

		// Call the contract method with the arguments
		e.callbackQueue.Enqueue(
			interfaces.NewContractMessage(
				contractId,
				method,
				args,
				e.contractId,
			),
		)
	}
}

func (e *Runtime) createContractEntry() func(caller *wasmtime.Caller, codeId int64, initArgsPtr int32) int32 {
	return func(caller *wasmtime.Caller, codeId int64, initArgsPtr int32) int32 {
		// TODO: consume fuel

		// Get the total contract amount as salt
		amount := e.repository.GetTotalContractAmount()
		salt := make([]byte, 8)
		binary.LittleEndian.PutUint64(salt, amount)

		contractId := generateContractId(e.state, uint64(codeId), salt)
		err := e.repository.CreateConctract(uint64(codeId), contractId)
		if err != nil {
			panic(err)
		}

		initArgs := readBytes(caller, initArgsPtr)

		e.callbackQueue.Enqueue(
			interfaces.NewContractMessage(
				contractId,
				"init",
				initArgs,
				e.contractId,
			),
		)

		// Return the contract ID as a pointer to the caller
		contractIdPtr := writeBytes(caller, []byte(contractId))
		return contractIdPtr
	}
}

func (e *Runtime) emitEventEntry() func(caller *wasmtime.Caller, eventPtr int32, dataPtr int32) {
	return func(caller *wasmtime.Caller, eventPtr int32, dataPtr int32) {
		// TODO: Consume fuel

		event := readUTF16EncodedString(readBytes(caller, eventPtr))
		data := readUTF16EncodedString(readBytes(caller, dataPtr))

		*e.resultEvents = append(*e.resultEvents, interfaces.ResultEvent{
			ContractId: e.contractId,
			Event:      event,
			Data:       data,
		})
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

func generateContractId(
	state []byte,
	codeId uint64,
	salt []byte,
) string {
	codeIdBz := make([]byte, 8)
	binary.LittleEndian.PutUint64(codeIdBz, codeId)
	contractId := sha256.Sum256(append(state, append(codeIdBz, salt...)...))
	return base58.Encode(contractId[:])
}
