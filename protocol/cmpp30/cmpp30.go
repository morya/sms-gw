package cmpp30

import (
	"github.com/morya/sms-gw/protocol/unified"
)

var c = &Coder{}

func init() {
	unified.RegisterProtocol("cmpp", c)
}
