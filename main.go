package main

import (
	"flag"

	"github.com/5dao/gdav/server"
)

func main() {
	configPath := flag.String("c", "config.toml", "-c config.toml,config file")
	flag.Parse()

	// locad config
	cfg, err := server.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	svr, err := server.NewServer(cfg)
	if err != nil {
		panic(err)
	}

	//log.Println(svr)
	go svr.Start()

	ch := make(chan int)
	<-ch
}
