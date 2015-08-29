package crxmake

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) TestLoadKeyFile(c *C) {
	b := NewBuilder()
	err := b.LoadKeyFile("fixtures/key.pem")
	c.Assert(err, IsNil)

	w := bytes.NewBuffer(nil)
	size, err := b.saveKeyFile(w)
	c.Assert(err, IsNil)
	c.Assert(size, Equals, 887)
	c.Assert(b.PrivateKey, Not(IsNil))
}

func (s *S) TestGenerateKeyIfNeeded(c *C) {
	b := NewBuilder()
	err := b.BuildZip("examples/app")
	c.Assert(err, IsNil)
	c.Assert(b.PrivateKey, Not(IsNil))
}

func (s *S) TestWriteToFile(c *C) {
	b := NewBuilder()
	err := b.LoadKeyFile("fixtures/key.pem")
	c.Assert(err, IsNil)

	err = b.BuildZip("examples/app")
	c.Assert(err, IsNil)

	w := bytes.NewBuffer(nil)
	err = b.WriteToFile(w)
	c.Assert(err, IsNil)

	hash := md5.Sum(w.Bytes())
	c.Assert(hex.EncodeToString(hash[:]), Equals, "e611edb20e227143f9bf7c68c28e1506")
}
