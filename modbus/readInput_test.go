package modbus

import (
	"gopkg.in/check.v1"
	"testing"
)

//Test is used to tap into golang testing framework
func Test(t *testing.T) { check.TestingT(t) }

type MySuite struct {
}

func (meSuit *MySuite) SetUpTest(c *check.C) {

}

func (meSuit *MySuite) TearDownTest(c *check.C) {

}

var _ = check.Suite(&MySuite{})

func (meSuit *MySuite) TestCore(lCheck *check.C) {
}

func (meSuit *MySuite) TestWhole(lCheck *check.C) {
}
