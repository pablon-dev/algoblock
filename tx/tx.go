package tx

import (
	"errors"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/algod/models"
)

var headers []*algod.Header
var client *algod.Client
var lastCheckedRound uint64

//Start needed parameters
func Start(round uint64, cli *algod.Client, head []*algod.Header) {
	lastCheckedRound = round
	client = cli
	headers = head
}

//GetNewTransactions from last checked
func GetNewTransactions(address string) ([]models.Transaction, error) {
	txParams, err := client.SuggestedParams()
	if err != nil {
		return nil, errors.New("Error getting suggested params")
	}
	currentRound := txParams.LastRound
	tlist, err := client.TransactionsByAddr(address, lastCheckedRound, currentRound, headers...)
	if err != nil {
		return nil, err
	}
	lastCheckedRound = currentRound
	return tlist.Transactions, nil
}

func GetNoteFields(txs ...models.Transaction) (result []string) {
	result = make([]string, len(txs))
	for i := range txs {
		result[i] = string(txs[i].Note)

	}
	return result
}

//WaitConf wait for confirmation
func WaitConf(txid string) {
	nodeStatus, err := client.Status()
	if err != nil {
		fmt.Println("Error getting status")
		return
	}
	lastRound := nodeStatus.LastRound
	for {
		pt, err := client.PendingTransactionInformation(txid)
		if err != nil {
			fmt.Println("waiting confirmations...")
			continue
		}
		if pt.ConfirmedRound > 0 {
			fmt.Printf("\n Transaction %s coinfirmed in round %d", pt.TxID, pt.ConfirmedRound)
			break
		}
		lastRound++
		client.StatusAfterBlock(lastRound)
	}
}
