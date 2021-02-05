package server

import (
	"fmt"
	"github.com/linsibolinhong/godis/cache"
	"github.com/linsibolinhong/godis/log"
	"github.com/linsibolinhong/godis/proto"
	"net"
)

type Server struct {
	port int
	c cache.Cache
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
		c:cache.NewSimpleCache(),
	}
}

func (s *Server) handler(c net.Conn) {
	proto_ := proto.NewRedisProto(c)
	for {
		cmd, err := proto_.ReadCommand()
		if err != nil {
			log.Error("read failed, err=%v", err)
			return
		}
		log.Info("cmd is %s", cmd.ToString())
		result, err := s.c.Command(cmd)
		if err != nil {
			log.Error("cmd error :%v", err)
			proto_.WriteError(err)
		} else {
			proto_.WriteResult(result)
		}
	}
}

func (s *Server) Run() error {
	server, err := net.Listen("tcp", ":" + fmt.Sprintf("%v", s.port))
	if err != nil {
		log.Error("%v", err)
		return err
	}

	log.Info("server is run on port:%v", s.port)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Error("conn error %v", err)
		} else {
			go s.handler(conn)
		}
	}

	return nil
}

