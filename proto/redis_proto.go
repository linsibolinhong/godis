package proto

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/linsibolinhong/godis/command"
	"github.com/linsibolinhong/godis/log"
)

const defaultBufSize = 4096

type redisProto struct {
	rw io.ReadWriter
	buf []byte
	bufIdx int
	bufLen int
}

func NewRedisProto(rw io.ReadWriter) Proto {
	return &redisProto{
		rw:rw,
		buf:make([]byte, defaultBufSize),
		bufIdx: 0,
		bufLen: 0,
	}
}

func (rp *redisProto) readArgNum() (int, error) {
	return rp.readNumber('*')
}

func (rp *redisProto) readParamLen() (int, error) {
	return rp.readNumber('$')
}

func (rp *redisProto) serial(s string) []byte {
	ret := make([]byte, 0, len(s)*2 + 2)
	ret = append(ret, '"')
	for _, c := range s {
		switch c {
		case '\\':
			ret = append(ret, '\\', '\\')
		case '\r':
			ret = append(ret, '\\', 'r')
		case '\n':
			ret = append(ret, '\\', 'n')
		case '"':
			ret = append(ret, '\\', '"')
		case '\t':
			ret = append(ret, '\\', 't')
		case '\b':
			ret = append(ret, '\\', 'b')
		default:
			ret = append(ret, byte(c))
		}
	}
	ret = append(ret, '"', '\r', '\n')
	return ret
}

func (rp *redisProto) readParam() ([]byte, error) {
	paramLen, err := rp.readParamLen()
	if err != nil {
		log.Error("read param len failed, err:%v", err)
		return nil, err
	}

	valbuf := make([]byte, paramLen)
	if rp.bufLen - rp.bufIdx >= paramLen {
		copy(valbuf, rp.buf[rp.bufIdx:rp.bufIdx + paramLen])
		rp.bufIdx += paramLen
		paramLen = 0
	} else {
		copy(valbuf, rp.buf[rp.bufIdx:])
		paramLen -= rp.bufLen - rp.bufIdx
		rp.bufIdx = rp.bufLen
	}

	for paramLen > 0 {
		idx := len(valbuf) - paramLen
		n, err := rp.rw.Read(valbuf[idx:])
		if err != nil {
			log.Error("read socket failed, err:%v", err)
			return nil, err
		}
		paramLen -= n;
	}

	err = rp.checkLineEnd()
	if err != nil {
		log.Error("check lined end failed, err:%v", err)
		return nil, err
	}

	return valbuf, nil
}

func (rp *redisProto) readFromSocket() error {
	for i := 0; i < 100; i++ {
		n, err := rp.rw.Read(rp.buf)
		if err != nil {
			log.Error("read socket failed : %v", err)
			return err
		}

		if n <= 0 {
			log.Info("read nothing, continue")
			time.Sleep(time.Millisecond)
			continue
		}

		rp.bufLen = n
		rp.bufIdx = 0
		return nil
	}

	return fmt.Errorf("read nothing")
}

func (rp *redisProto) printBuf() {
	log.Debug("%s", string(rp.buf[rp.bufIdx:]))
}

func (rp *redisProto) readNumber(flag byte) (int, error) {
	checkFlags := false
	num := 0
	maxNum := 1024 * 1024 * 1024

	for {
		if rp.bufIdx >= rp.bufLen {
			if err := rp.readFromSocket(); err != nil {
				log.Error("read from socket failed, err:%v", err)
				return -1, err
			}
		}

		if !checkFlags {
			if rp.buf[rp.bufIdx] != flag {
				if flag == '*' {
					err := rp.checkLineEnd()
					if err != nil {
						log.Error("invalid flag %v != %v", flag,rp.buf[rp.bufIdx])
						return -1, err
					}
					return 0, nil
				}
				log.Error("invalid flag %v != %v", string(flag),string(rp.buf[rp.bufIdx]))
				rp.printBuf()
				return -1, fmt.Errorf("invalid flag")
			} else {
				checkFlags = true
				rp.bufIdx ++
				continue
			}
		}

		for ; rp.bufIdx < rp.bufLen; {
			c := rp.buf[rp.bufIdx]
			if c >= '0' && c <= '9' {
				num = num * 10 + int(c - '0')
				if num > maxNum {
					log.Error("too large number")
					return -1, fmt.Errorf("too large number")
				}
				rp.bufIdx++
				continue
			}

			err := rp.checkLineEnd()
			if err != nil {
				log.Error("check line end failed", err)
				return -1, fmt.Errorf("check line end failed")
			}

			return num, nil
		}
	}
}

func (rp *redisProto) checkLineEnd() error {
	checked := 0
	for {
		if rp.bufIdx >= rp.bufLen {
			err := rp.readFromSocket()
			if err != nil {
				log.Error("read socket failed, err:%v", err)
				return fmt.Errorf("read socket failed")
			}
		}

		switch checked {
		case 0:
			if rp.buf[rp.bufIdx] != '\r' {
				return fmt.Errorf("check end failed")
			}
			checked = 1
			rp.bufIdx++
		case 1:
			if rp.buf[rp.bufIdx] != '\n' {
				return fmt.Errorf("check end failed")
			}
			rp.bufIdx += 1
			return nil
		default:
			return fmt.Errorf("nevever reached")
		}
	}

}

func (rp *redisProto) ReadCommand() (*command.Command, error) {
	argNum, err := rp.readArgNum()
	cmd := command.NewCommand()

	if err != nil {
		log.Error("read argnum failed, err:%v", err)
		return nil, err
	}

	log.Info("arg num is %d", argNum)
	for i := 0; i < argNum; i++ {
		param, err := rp.readParam()
		if err != nil {
			log.Error("read param failed, err:%v", err)
			return nil, err
		}
		if i == 0 {
			cmd.Cmd = command.Method(strings.ToLower(string(param)))
		} else {
			cmd.AppendParam(string(param))
		}
	}
	return cmd, nil
}

func (rp *redisProto) WriteResult(result *command.Result) error {
	if result == nil || len(result.Ret) == 0 {
		_, e := rp.rw.Write([]byte("+OK\r\n"))
		return e
	}

	if len(result.Ret) == 1 {
		_, e := rp.rw.Write(append([]byte("+"), rp.serial(result.Ret[0])...))
		return e
	}

	ret := []byte(fmt.Sprintf("*%v\r\n", len(result.Ret)))
	for _, p := range result.Ret {
		ret = append(ret, []byte(fmt.Sprintf("$%d\r\n", len(p)))...)
		ret = append(ret, rp.serial(p)...)

	}
	_, e := rp.rw.Write([]byte(ret))
	return e
}

func (rp *redisProto) WriteError(err error) error {
	ret := "-" + err.Error() + "\r\n"
	_, e := rp.rw.Write([]byte(ret))
	return e
}