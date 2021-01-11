package main

import (
	"fmt"
	"github.com/linsibolinhong/godis/server"
)

func main() {
	ser := server.NewServer(6789)
	fmt.Println(ser.Run())
}
