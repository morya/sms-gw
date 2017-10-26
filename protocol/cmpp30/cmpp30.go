package cmpp30

import (
	"github.com/morya/sms/protocol/unified"
)

var c = &Coder{}

func init() {
	unified.RegisterProtocol("cmpp", c)
}
