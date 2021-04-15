package databases

import (
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

	_, err = tarConn.Ping()
	if err != nil {
		return nil, err
	}

	return &Tarantool{
		tarantoolDatabase: tarConn,
	}, nil
}

func (t *Tarantool) GetDatabase() *tarantool.Connection {
	return t.tarantoolDatabase
}
