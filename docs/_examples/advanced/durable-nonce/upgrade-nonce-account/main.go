package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yosefl20/solana-go-sdk/client"
	"github.com/yosefl20/solana-go-sdk/common"
	"github.com/yosefl20/solana-go-sdk/program/system"
	"github.com/yosefl20/solana-go-sdk/rpc"
	"github.com/yosefl20/solana-go-sdk/types"
)

// FUarP2p5EnxD66vVDL4PWRoWMzA56ZVHG24hpEDFShEz
var feePayer, _ = types.AccountFromBase58("4TMFNY9ntAn3CHzguSAvDNLPRoQTaK3sWbQQXdDXaE6KWRBLufGL6PJdsD2koiEe3gGmMdRK3aAw7sikGNksHJrN")

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	nonceAccountPubkey := common.PublicKeyFromString("DJyNpXgggw1WGgjTVzFsNjb3fuQZVMqhoakvSBfX9LYx")

	// recent blockhash
	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{feePayer},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				system.UpgradeNonceAccount(system.UpgradeNonceAccountParam{
					NonceAccountPubkey: nonceAccountPubkey,
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a transaction, err: %v", err)
	}

	sig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	fmt.Println("txhash", sig)
}
