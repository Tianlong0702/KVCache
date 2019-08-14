package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"KVCache/cache"
	"KVCache/server"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		log.Fatal("Please provide a port number!")
		return
	}

	Port := ":" + arguments[1]

	cacheType := cache.LRU
	cacheCapacity := 100

	if len(arguments) > 2 {
		cacheType = arguments[2]
	}
	var err error
	if len(arguments) > 3 {
		cacheCapacity, err = strconv.Atoi(arguments[3])
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	// initialize local cache engine
	localCache := cache.NewCache(cacheType)
	localCache.Init(cacheCapacity)

	// create a new key-value cache server
	kvCaheServer := server.New(Port, localCache)

	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		// shut down
		<-osSignal
		kvCaheServer.Stop()
		os.Exit(1)
	}()

	// start to receive incoming request from network
	kvCaheServer.Start()
}
