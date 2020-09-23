package asset

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/pablon-dev/AlgoBlock/tx"
	"github.com/pablon-dev/AlgoBlock/util"
)

//Create the default block assest
func Create(client *algod.Client, headers []*algod.Header, mnems []string) error {
	txParams, err := client.SuggestedParams()
	if err != nil {
		return errors.New("Error getting suggested params")
	}
	util.Pretty(txParams)
	if len(mnems) < 2 {
		return errors.New("Invalid mnemonics")
	}
	pivk1, err := mnemonic.ToPrivateKey(mnems[0])
	if err != nil {
		return err
	}
	pivk2, err := mnemonic.ToPrivateKey(mnems[1])
	if err != nil {
		return err
	}
	pub1, err := crypto.GenerateAddressFromSK([]byte(pivk1))
	if err != nil {
		return err
	}
	pub2, err := crypto.GenerateAddressFromSK([]byte(pivk2))
	if err != nil {
		return err
	}
	addr1 := pub1.String()
	addr2 := pub2.String()

	fee := txParams.Fee
	firstRound := txParams.LastRound
	lastRound := txParams.LastRound + 1000
	genhash := base64.StdEncoding.EncodeToString(txParams.GenesisHash)
	genID := txParams.GenesisID
	creator := addr1
	assetName := "BLQ Vote"
	unitName := "BLQV"
	assetURL := "https://peakd.com/@cubanblock"
	metaHash := "d41d8cd98f00b204e9800998ecf8427e"
	frozen := false
	decimals := uint32(2)
	supply := uint64(3000000)
	manager := addr2
	note := []byte("Inicio!!")
	//TX
	txn, err := transaction.MakeAssetCreateTxn(creator, fee, firstRound, lastRound, note, genID, genhash, supply, decimals, frozen, manager, "", "", "", unitName, assetName, assetURL, metaHash)
	if err != nil {
		return errors.New("Failed creating TX")
	}
	fmt.Printf("Asset created: %s\n", assetName)
	txid, stx, err := crypto.SignTransaction(pivk1, txn)
	if err != nil {
		return errors.New("Failed signing TX")
	}
	fmt.Printf("txid: %s", txid)

	resp, err := client.SendRawTransaction(stx)
	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to broadcast TX")
	}
	tx.WaitConf(resp.TxID)

	//Info
	act, err := client.AccountInformation(addr1, headers...)
	if err != nil {
		return err
	}
	assetID := uint64(0)
	for i := range act.AssetParams {
		if i > assetID {
			assetID = i
		}
	}
	fmt.Println("AssetID: ", assetID)
	assetinfo, err := client.AssetInformation(assetID, headers...)
	if err != nil {
		return err
	}
	util.Pretty(assetinfo)
	return nil
}
