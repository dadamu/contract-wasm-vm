package store

import (
	"testing"

	"github.com/cosmos/iavl"
	dbm "github.com/cosmos/iavl/db"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	tree  *iavl.MutableTree
	store *Store
	cache *CacheKVStore
}

func (s *TestSuite) SetupTest() {
	s.tree = iavl.NewMutableTree(dbm.NewMemDB(), 100, false, iavl.NewNopLogger())
	s.store = NewStore(s.tree)
	s.cache = s.store.GetCached()
}

func TestRuntimeTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// -----------------------------------------------------------------------------

func (s *TestSuite) TestUncommit() {
	s.cache.set([]byte("key1"), []byte("value1"))

	// Check that the value is not in the tree since it is only in the repository
	valueInTree, err := s.tree.Get([]byte("key1"))
	s.Require().NoError(err)
	s.Require().Nil(valueInTree)
}

func (s *TestSuite) TestCommit() {
	s.cache.set([]byte("key1"), []byte("value1"))
	s.cache.Commit()

	// Check that the value is in the tree after commit
	valueInTree, err := s.tree.Get([]byte("key1"))
	s.Require().NoError(err)
	s.Require().Equal([]byte("value1"), valueInTree)
}

func (s *TestSuite) TestSaveAndLoadVersion() {
	s.cache.set([]byte("key1"), []byte("value1"))
	s.cache.Commit()

	expectedVersion := s.tree.WorkingVersion()
	_, err := s.store.SaveVersionWithId(1)
	s.Require().NoError(err)

	// Check that the version is saved correctly
	version, err := s.store.GetVersionById(1)
	s.Require().NoError(err)
	s.Require().Equal(uint64(expectedVersion), version)
}
