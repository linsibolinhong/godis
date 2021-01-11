package server

import (
	"fmt"
	"github.com/linsibolinhong/godis/log"
	"github.com/linsibolinhong/godis/proto"
	"net"
)

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{port: port}
}

func (s *Server) handler(c net.Conn) {
	proto_ := proto.NewRedisProto(c)
	for {
		_, err := proto_.ReadCommand()
		if err != nil {
			log.Error("read failed, err=%v", err)
			return
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

