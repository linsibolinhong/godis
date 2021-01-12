package main

import (
	"fmt"
	"github.com/linsibolinhong/godis/server"
)

func main() {
	ser := server.NewServer(6379)
	fmt.Println(ser.Run())
}
