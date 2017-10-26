package sgip12

import (
	"io"

	"github.com/morya/sms/protocol/unified"
)

var c = &Coder{}

func init() {
	unified.RegisterProtocol("sgip", c)
}

type Coder struct {
}

func (c *Coder) NextMsg(r io.Reader) (data []byte, err error) {
	return
}
func (c *Coder) Decode([]byte) (interface{}, error) {
	return nil, nil

}
func (c *Coder) Encode(v interface{}) (data []byte, err error) {
	return
}
