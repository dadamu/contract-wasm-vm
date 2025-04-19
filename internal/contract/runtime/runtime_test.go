package runtime

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/bytecodealliance/wasmtime-go/v31"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	callbackqueue "github.com/dadamu/contract-wasmvm/internal/contract/callback-queue"
	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces"
	"github.com/dadamu/contract-wasmvm/internal/contract/interfaces/testutil"
)

type RuntimeTestSuite struct {
	suite.Suite
	queue      *callbackqueue.CallbackQueue
	events     *[]interfaces.ResultEvent
	repository *testutil.MockIContractRepository
	runtime    *Runtime
}

func (suite *RuntimeTestSuite) SetupTest() {
	gomockCtrl := gomock.NewController(suite.T())

	// Read wasm from file
	wasmFile, err := os.ReadFile("testdata/test.wasm")
	if err != nil {
		suite.T().Fatalf("failed to read module: %v", err)
	}

	suite.repository = testutil.NewMockIContractRepository(gomockCtrl)
	suite.queue = callbackqueue.NewCallbackQueue()
	suite.events = &[]interfaces.ResultEvent{}

	config := wasmtime.NewConfig()
	config.SetConsumeFuel(true)
	engine := wasmtime.NewEngineWithConfig(config)
	module, err := wasmtime.NewModule(engine, wasmFile)
	if err != nil {
		suite.T().Fatalf("failed to create module: %v", err)
	}

	suite.runtime = NewRuntimeFromModule(
		engine,
		suite.queue,
		suite.events,
		suite.repository,
		module,

		[]byte("state"),
		"contractId",
		20_000,
	)
}

func TestRuntimeTestSuite(t *testing.T) {
	suite.Run(t, new(RuntimeTestSuite))
}

// -----------------------------------------------------------------------------

func (s *RuntimeTestSuite) TestInfiniteLoop() {
	_, err := s.runtime.Run(interfaces.NewContractMessage("contractId", "infiniteLoop", []byte{}, "sender"))
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "wasm trap: all fuel consumed by WebAssembly")
}

func (s *RuntimeTestSuite) TestSaveAndLoad() {
	// Setup the mock repository
	loadedValue := []byte{1, 0, 0, 0} // Initial value: 1
	savedValue := []byte{2, 0, 0, 0}  // Expected value after increment: 2

	s.repository.EXPECT().LoadEntity("contractId", "test").Return(loadedValue)
	s.repository.EXPECT().SaveEntity("contractId", "test", savedValue)

	_, err := s.runtime.Run(interfaces.NewContractMessage("contractId", "addOne", []byte{}, "sender"))
	s.Require().NoError(err)
}

func (s *RuntimeTestSuite) TestGasRemaining() {
	// Setup the mock repository
	loadedValue := []byte{1, 0, 0, 0} // Initial value: 1
	savedValue := []byte{2, 0, 0, 0}  // Expected value after increment: 2

	s.repository.EXPECT().LoadEntity("contractId", "test").Return(loadedValue)
	s.repository.EXPECT().SaveEntity("contractId", "test", savedValue)

	remaining, err := s.runtime.Run(interfaces.NewContractMessage("contractId", "addOne", []byte{}, "sender"))
	s.Require().NoError(err)
	s.Require().Equal(uint64(8_296), remaining)
}

func (s *RuntimeTestSuite) TestCrash() {
	_, err := s.runtime.Run(interfaces.NewContractMessage("contractId", "crash", []byte{}, "sender"))
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "WASM called abort msg:")
}

func (s *RuntimeTestSuite) TestContractCall() {
	_, err := s.runtime.Run(interfaces.NewContractMessage("contractId", "callback", []byte{}, "sender"))
	s.Require().NoError(err)

	msg, found := s.queue.Dequeue()
	s.Require().True(found)
	s.Require().Equal("contract", msg.Contract)
	s.Require().Equal("method", msg.Method)
	s.Require().Equal([]byte("args"), msg.Args)
}

func (s *RuntimeTestSuite) TestCreateContract() {
	salt := make([]byte, 8)
	binary.LittleEndian.PutUint64(salt, 1)

	contractId := generateContractId(s.runtime.state, uint64(1), salt)

	s.repository.EXPECT().GetTotalContractAmount().Return(uint64(1))
	s.repository.EXPECT().CreateConctract(uint64(1), contractId)

	_, err := s.runtime.Run(interfaces.NewContractMessage("contractId", "createContract", []byte{}, "sender"))
	s.Require().NoError(err)

	msg, found := s.queue.Dequeue()
	s.Require().True(found)
	s.Require().Equal(contractId, msg.Contract)
	s.Require().Equal("init", msg.Method)
	s.Require().Equal([]byte("args"), msg.Args)
	s.Require().Equal(s.runtime.contractId, msg.Sender)
}

func (s *RuntimeTestSuite) TestEmitEvent() {
	_, err := s.runtime.Run(interfaces.NewContractMessage("contractId", "emitEvent", []byte{}, "sender"))
	s.Require().NoError(err)

	event := (*s.events)[0]

	s.Require().Equal("contractId", event.ContractId)
	s.Require().Equal("event", event.Event)
	s.Require().Equal("data", event.Data)
}
