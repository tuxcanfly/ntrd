package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lightninglabs/neutrino"
	"github.com/roasbeef/btcd/chaincfg"
	"github.com/roasbeef/btcwallet/walletdb"
	_ "github.com/roasbeef/btcwallet/walletdb/bdb"
)

func main() {
	db, err := walletdb.Create("bdb", "/tmp/neutrino/wallet.db")
	if err != nil {
		log.Fatalf("unable to create db: %v", err)
	}
	svc, err := neutrino.NewChainService(neutrino.Config{
		DataDir:      "/tmp/neutrino",
		Database:     db,
		Namespace:    []byte("neutrino"),
		ChainParams:  chaincfg.SimNetParams,
		ConnectPeers: []string{"127.0.0.1:12555"},
	})
	if err != nil {
		log.Fatalf("unable to create neutrino: %v", err)
	}
	svc.Start()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		svc.Stop()
		os.Exit(1)
	}()

	for {
		fmt.Println("...")
		time.Sleep(10 * time.Second)
	}
}
