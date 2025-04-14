package runtime

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go/v31"

	callbackqueue "github.com/dadamu/contract-wasmvm/internal/contract/callback-queue"
	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces"
)

type Runtime struct {
	callbackQueue *callbackqueue.CallbackQueue

	engine   *wasmtime.Engine
	store    *wasmtime.Store
	instance *wasmtime.Instance

	contractId string
	repository interfaces.IContractRepository
}

func NewRuntimeFromModule(
	callbackQueue *callbackqueue.CallbackQueue,
	engine *wasmtime.Engine,
	contractId string,
	repository interfaces.IContractRepository,
	module *wasmtime.Module,
	gasLimit uint64,
) *Runtime {
	runtime := &Runtime{
		callbackQueue: callbackQueue,
		engine:        engine,
		contractId:    contractId,
		repository:    repository,
	}

	instance := runtime.newInstanceFromModule(module, gasLimit)
	if instance == nil {
		panic("failed to create instance")
	}

	runtime.instance = instance
	return runtime
}

func (e *Runtime) newInstanceFromModule(module *wasmtime.Module, gasLimit uint64) *wasmtime.Instance {
	// Create a new isolated store for the contract runtime
	store := wasmtime.NewStore(e.engine)
	store.SetFuel(gasLimit)
	e.store = store
	linker := e.prepareLinker()

	instance, err := linker.Instantiate(store, module)
	if err != nil {
		panic(err)
	}
	return instance
}

func (e *Runtime) Run(state []byte, msg interfaces.ContractMessage) (remainingGas uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	run := e.instance.GetFunc(e.store, msg.Method)
	if run == nil {
		return 0, fmt.Errorf("function %s not found in instance", msg.Method)
	}

	senderPtr := e.writeBytesByInstance([]byte(msg.Sender))
	statePtr := e.writeBytesByInstance(state)
	argsPtr := e.writeBytesByInstance(msg.Args)

	// Call the run function with the pointer to the golobal state and args
	_, err = run.Call(e.store, statePtr, senderPtr, argsPtr)
	if err != nil {
		return 0, err
	}

	// Get the remaining fuel
	// Ignore error that fuel is not configured for the store
	remainingGas, _ = e.store.GetFuel()

	return remainingGas, nil
}

func (e *Runtime) writeBytesByInstance(message []byte) int32 {
	malloc := e.instance.GetExport(e.store, "__new").Func()
	if malloc == nil {
		panic("__new function not found")
	}

	resultPtr, err := malloc.Call(e.store, int32(len(message)), 1)
	if err != nil {
		panic(err)
	}
	offset := resultPtr.(int32)
	memory := e.instance.GetExport(e.store, "memory").Memory().UnsafeData(e.store)
	copy(memory[offset:offset+int32(len(message))], message)
	return offset
}
