package main

import (
	node "bitcask_master_slave/node/kitex_gen/node/nodeservice"
	"bitcask_master_slave/pkg/consts"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	addr, err := net.ResolveTCPAddr(consts.TCP, consts.NodeAddr)
	if err != nil {
		panic(err)
	}

	svr := node.NewServer(
		new(NodeServiceImpl),
		server.WithServiceAddr(addr),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
