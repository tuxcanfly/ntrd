package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/btcsuite/btclog"
	"github.com/lightninglabs/neutrino"
	"github.com/roasbeef/btcd/chaincfg"
	"github.com/roasbeef/btcwallet/walletdb"
	_ "github.com/roasbeef/btcwallet/walletdb/bdb"
)

func main() {
	logger := btclog.NewBackend(os.Stdout).Logger("NTRD")
	logger.SetLevel(btclog.LevelDebug)
	neutrino.UseLogger(logger)
	db, err := walletdb.Create("bdb", "/tmp/neutrino/wallet.db")
	if err != nil {
		log.Fatalf("unable to create db: %v", err)
		return
	}
	svc, err := neutrino.NewChainService(neutrino.Config{
		DataDir:      "/tmp/neutrino",
		Database:     db,
		Namespace:    []byte("neutrino"),
		ChainParams:  chaincfg.TestNet3Params,
		ConnectPeers: []string{"btcd0.lightning.computer:18333"},
	})
	if err != nil {
		logger.Errorf("unable to create neutrino: %v", err)
		return
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
