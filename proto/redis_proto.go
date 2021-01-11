package proto

import (
	"github.com/linsibolinhong/godis/command"
	"github.com/linsibolinhong/godis/log"
	"io"
)

const defaultBufSize = 4096

type redisProto struct {
	rw io.ReadWriter
	buf []byte
}

func NewRedisProto(rw io.ReadWriter) Proto {
	return &redisProto{
		rw:rw,
		buf:make([]byte, defaultBufSize),
	}
}

func (rp *redisProto) ReadCommand() (*command.Command, error) {
	for {
		n, err := rp.rw.Read(rp.buf)
		if err != nil {
			log.Error("read socket failed, err:%v", err)
			return nil, err
		}
		log.Error("%v, %v", n, string(rp.buf[:n]))
	}

	return nil, nil
}

func (rp *redisProto) WriteResult(result *command.Result) error {
	return nil
}