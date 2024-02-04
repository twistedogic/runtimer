package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/twistedogic/runtimer/internal"
)

var (
	portVar   int
	hostVar   string
	configVar string
)

func init() {
	flag.IntVar(&portVar, "port", 3000, "port to listen")
	flag.StringVar(&hostVar, "hostname", "127.0.0.1", "hostname to listen")
	flag.StringVar(&configVar, "config", "", "path of the config file")
}

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", hostVar, portVar)
	ctx, cancel := context.WithCancel(context.Background())
	runner, err := internal.NewRunner(configVar)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	c := make(chan os.Signal)
	signal.Notify(
		c,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2,
	)
	go func() {
		<-c
		cancel()
	}()
	if err := runner.Start(ctx, addr); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
