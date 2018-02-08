package rpc

import (
	"apiservice/download"
	"apiservice/g"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
	// "happy-hbs/modules/hbs/download"
	// "happy-hbs/modules/hbs/g"
)

type Hbs int
type Agent int
type Plugin int

func Start() {
	// Init Download File Set
	download.NewFileSet()

	addr := g.Config().Listen

	server := rpc.NewServer()
	// server.Register(new(filter.Filter))
	server.Register(new(Agent))
	server.Register(new(Hbs))
	server.Register(new(Plugin))

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalln("listen error:", e)
	} else {
		log.Println("listening", addr)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("listener accept fail:", err)
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
