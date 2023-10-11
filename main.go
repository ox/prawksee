package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

var configFilepath string

func init() {
	flag.StringVar(&configFilepath, "config", "", "Path to config file")
	flag.Parse()

	if configFilepath == "" {
		fmt.Printf("missing required --config flag\n")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	fmt.Printf("config at: %s\n", configFilepath)

	config, err := ParseConfig(configFilepath)
	if err != nil {
		fmt.Printf("could not parse config: %s", err)
		os.Exit(1)
	}

	listeners := make(map[string]net.Listener)
	for name, c := range config.Servers {
		listener, err := net.Listen(c.Network, fmt.Sprintf(":%d", c.Bind))
		if err != nil {
			log.Printf("could not listen on %d for server %s: %s", c.Bind, name, err)
			os.Exit(1)
		}

		listeners[name] = listener
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	go func() {
		select {
		case <-sigChan:
			for _, listener := range listeners {
				listener.Close()
			}
		}
	}()

	wg := sync.WaitGroup{}
	for serverName, serverInstance := range listeners {
		serverConfig := config.Servers[serverName]
		wg.Add(1)
		go func(n string, c Server, s net.Listener) {
			log.Printf("[%s] routing :%d to %s", n, c.Bind, c.Address)
			for {
				conn, err := s.Accept()
				if err != nil {
					if !errors.Is(err, net.ErrClosed) {
						log.Panicf("[%s] error accepting connection: %s", n, err)
					}
					break
				}

				dest, err := net.Dial(c.Network, c.Address)
				if err != nil {
					log.Printf("[%s] could not dial upstream: %s", n, err)
					conn.Close()
					continue
				}

				go transfer(conn, dest)
			}
			log.Printf("[%s] shutting down", n)
			wg.Done()
		}(serverName, serverConfig, serverInstance)
	}

	wg.Wait()
}

func transfer(conn net.Conn, dest net.Conn) {
	go func() {
		defer func() {
			conn.Close()
		}()

		io.Copy(conn, dest)
	}()

	go func() {
		defer func() {
			dest.Close()
		}()

		io.Copy(dest, conn)
	}()
}
