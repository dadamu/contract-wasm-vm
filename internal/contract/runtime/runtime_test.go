package runtime

import (
	"os"
	"testing"

	"github.com/bytecodealliance/wasmtime-go/v31"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	callbackqueue "github.com/dadamu/contract-wasmvm/internal/contract/callback-queue"
	repositorytestutil "github.com/dadamu/contract-wasmvm/internal/contract/repository/testutil"
)

type RuntimeTestSuite struct {
	suite.Suite
	queue      *callbackqueue.CallbackQueue
	repository *repositorytestutil.MockIRepository
	runtime    *Runtime
}

func (suite *RuntimeTestSuite) SetupTest() {
	gomockCtrl := gomock.NewController(suite.T())

	// Read wasm from file
	wasmFile, err := os.ReadFile("testdata/release.wasm")
	if err != nil {
		suite.T().Fatalf("failed to read module: %v", err)
	}

	suite.repository = repositorytestutil.NewMockIRepository(gomockCtrl)
	suite.queue = callbackqueue.NewCallbackQueue()

	config := wasmtime.NewConfig()
	config.SetConsumeFuel(true)
	engine := wasmtime.NewEngineWithConfig(config)
	modele, err := wasmtime.NewModule(engine, wasmFile)
	if err != nil {
		suite.T().Fatalf("failed to create module: %v", err)
	}

	runtime := NewRuntimeFromModule(
		suite.queue,
		engine,
		"test",
		suite.repository,
		modele,
		20_000,
	)

	suite.runtime = runtime
}

func (s *RuntimeTestSuite) TestInfiniteLoop() {
	_, err := s.runtime.Run("infiniteLoop", []byte{}, []byte{})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "wasm trap: all fuel consumed by WebAssembly")
}

func (s *RuntimeTestSuite) TestSaveAndLoad() {
	// Setup the mock repository
	loadedValue := []byte{1, 0, 0, 0} // Initial value: 1
	savedValue := []byte{2, 0, 0, 0}  // Expected value after increment: 2

	s.repository.EXPECT().LoadEntity("test", "test").Return(loadedValue, nil)
	s.repository.EXPECT().SaveEntity("test", "test", savedValue).Return(nil)

	_, err := s.runtime.Run("addOne", []byte{}, []byte{})
	s.Require().NoError(err)
}

func (s *RuntimeTestSuite) TestGasRemaining() {
	// Setup the mock repository
	loadedValue := []byte{1, 0, 0, 0} // Initial value: 1
	savedValue := []byte{2, 0, 0, 0}  // Expected value after increment: 2

	s.repository.EXPECT().LoadEntity("test", "test").Return(loadedValue, nil)
	s.repository.EXPECT().SaveEntity("test", "test", savedValue).Return(nil)

	remaining, err := s.runtime.Run("addOne", []byte{}, []byte{})
	s.Require().NoError(err)
	s.Require().Equal(uint64(8_045), remaining)
}

func (s *RuntimeTestSuite) TestCrash() {
	_, err := s.runtime.Run("crash", []byte{}, []byte{})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "WASM called abort msg:")
}

func (s *RuntimeTestSuite) TestContractCall() {
	_, err := s.runtime.Run("callback", []byte{}, []byte{})
	s.Require().NoError(err)

	msg, found := s.queue.Dequeue()
	s.Require().True(found)
	s.Require().Equal("contract", msg.Contract)
	s.Require().Equal("method", msg.Method)
	s.Require().Equal([]byte("args"), msg.Args)
}

func TestRuntimeTestSuite(t *testing.T) {
	suite.Run(t, new(RuntimeTestSuite))
}
