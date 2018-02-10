package main

import (
	"apiservice/g"
	"apiservice/http"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	daemon := flag.Bool("d", false, "run in daemon mode")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *daemon {
		fmt.Println("runing in daemon.")
	}

	g.ParseConfig(*cfg)

	// db.Init()
	// cache.Init()
	// g.InitRedisConnPool()
	// sender.Start()

	// go cache.DeleteStaleAgents()

	go http.Start()
	// go rpc.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		// fmt.Println()
		// db.DB.Close()
		// g.RedisConnPool.Close()
		os.Exit(0)
	}()

	select {}
}
