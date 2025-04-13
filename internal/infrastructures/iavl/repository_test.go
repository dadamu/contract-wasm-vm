package iavl

import (
	"testing"

	"github.com/cosmos/iavl"
	dbm "github.com/cosmos/iavl/db"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	tree *iavl.MutableTree
}

func (s *TestSuite) SetupTest() {
	s.tree = iavl.NewMutableTree(dbm.NewMemDB(), 100, false, iavl.NewNopLogger())
}

func TestRuntimeTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// -----------------------------------------------------------------------------

func (s *TestSuite) TestSaveAndLoad() {
	repo := NewIAVLRepository("test", s.tree)
	err := repo.SaveEntity("test", "test", []byte{1, 0, 0, 0})
	s.Require().NoError(err)

	value, err := repo.LoadEntity("test", "test")
	s.Require().NoError(err)
	s.Require().Equal([]byte{1, 0, 0, 0}, value)
}

func (s *TestSuite) TestLoadEntityNotFound() {
	repo := NewIAVLRepository("test", s.tree)
	value, err := repo.LoadEntity("test", "notfound")
	s.Require().NoError(err)
	s.Require().Nil(value)
}
