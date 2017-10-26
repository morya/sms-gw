package unified

import (
	"fmt"
	"io"
	"sync"
)

type Coder interface {
	NextMsg(r io.Reader) ([]byte, error)
	Decode([]byte) (interface{}, error)
	Encode(v interface{}) ([]byte, error)
}

var globalProtolRegister = make(map[string]Coder)
var lock sync.Mutex

func RegisterProtocol(name string, c Coder) error {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := globalProtolRegister[name]; ok {
		return fmt.Errorf("we already have protocol named %s", name)
	}

	globalProtolRegister[name] = c
	return nil
}

func GetCoder(name string) Coder {
	lock.Lock()
	defer lock.Unlock()

	c := globalProtolRegister[name]
	return c
}

func ListProtocols() []string {
	var names []string

	lock.Lock()
	defer lock.Unlock()

	for name := range globalProtolRegister {
		names = append(names, name)
	}

	return names
}
