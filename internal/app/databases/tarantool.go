package databases

import (
	"fmt"

	"github.com/tarantool/go-tarantool"
)

type Tarantool struct {
	tarantoolDatabase *tarantool.Connection
}

func NewTarantool(user string, pass string, addr string) (*Tarantool, error) {
	opts := tarantool.Opts{
		User: user,
		Pass: pass,
	}

	tarConn, err := tarantool.Connect(addr, opts)
	if err != nil {
		return nil, err
	}

	resp, err := tarConn.Ping()
	fmt.Println(resp, err)
	return &Tarantool{
		tarantoolDatabase: tarConn,
	}, nil
}

func (t *Tarantool) GetDatabase() *tarantool.Connection {
	return t.tarantoolDatabase
}
