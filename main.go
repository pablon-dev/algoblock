package main

import (
	"crypto/ed25519"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/pablon-dev/AlgoBlock/accounts"
	"github.com/pablon-dev/AlgoBlock/hive"
	"github.com/pablon-dev/AlgoBlock/tx"
)

var (
	algoRecieverAddr string
	algoApiKey       string
)

const algodServer string = "https://mainnet-algorand.api.purestake.io/ps1"
const assetRecID uint64 = 36916093

var mnems, pubks []string
var pivks []ed25519.PrivateKey

func init() {
	algoRecieverAddr = os.Getenv("algo_reciever")
	algoApiKey = os.Getenv("algoApiKey")
	if err := loadmnemonics(); err != nil {
		panic(err)
	}
	if err := loadkeys(); err != nil {
		panic(err)
	}

}
func loadmnemonics() error {
	var err error
	mnems, err = accounts.ReadMnemonics()
	if err != nil {
		return err
	}
	return nil
}
func loadkeys() error {
	pivks = make([]ed25519.PrivateKey, len(mnems))
	pubks = make([]string, len(mnems))
	for i, m := range mnems {
		pv, err := mnemonic.ToPrivateKey(m)
		if err != nil {
			return err
		}
		pivks[i] = pv
		pk := pv.Public()
		var a types.Address
		cpk := pk.(ed25519.PublicKey)
		copy(a[:], cpk[:])
		pubks[i] = a.String()
		fmt.Println("Load Key:", pubks[i])
	}
	return nil

}
func main() {
	var headers []*algod.Header
	headers = append(headers, &algod.Header{"X-API-Key", algoApiKey})
	client, err := algod.MakeClientWithHeaders(algodServer, "", headers)
	if err != nil {
		fmt.Println("failed starting algod")
		return
	}
	txParams, err := client.SuggestedParams()
	if err != nil {
		fmt.Println("Error getting params")
		return
	}
	tx.Start(txParams.LastRound, &client, headers)
	fmt.Printf("Started at block: %d, listening for transaction to AssetID: %d", txParams.LastRound, assetRecID)

	go func() {
		for now := range time.Tick(time.Second * 50) {
			fmt.Println("Checked at:", now)
			checkForTransactions()
		}
	}()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}
func checkForTransactions() {
	var txs []models.Transaction
	var err error
	txs, err = tx.GetNewTransactions(algoRecieverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(txs) > 0 {
		var confTxS []string
		var powers []uint64
		for i := range txs {
			if string(txs[i].Type) != "axfer" || (*txs[i].AssetTransfer).AssetID != assetRecID {
				continue
			}
			confTxS = append(confTxS, string(txs[i].Note))
			powers = append(powers, (*txs[i].AssetTransfer).Amount)
		}
		hive.CastVote(confTxS, powers)
	}
}

//Pretty dibujitos
