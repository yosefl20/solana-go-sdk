package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yosefl20/solana-go-sdk/client"
	"github.com/yosefl20/solana-go-sdk/common"
	"github.com/yosefl20/solana-go-sdk/pkg/pointer"
	"github.com/yosefl20/solana-go-sdk/program/associated_token_account"
	"github.com/yosefl20/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/yosefl20/solana-go-sdk/program/system"
	"github.com/yosefl20/solana-go-sdk/program/token"
	"github.com/yosefl20/solana-go-sdk/rpc"
	"github.com/yosefl20/solana-go-sdk/types"
)

// FUarP2p5EnxD66vVDL4PWRoWMzA56ZVHG24hpEDFShEz
var feePayer, _ = types.AccountFromBase58("4TMFNY9ntAn3CHzguSAvDNLPRoQTaK3sWbQQXdDXaE6KWRBLufGL6PJdsD2koiEe3gGmMdRK3aAw7sikGNksHJrN")

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	mint := types.NewAccount()
	fmt.Printf("NFT: %v\n", mint.PublicKey.ToBase58())

	ata, _, err := common.FindAssociatedTokenAddress(feePayer.PublicKey, mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid ata, err: %v", err)
	}

	tokenMetadataPubkey, err := token_metadata.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid token metadata, err: %v", err)

	}

	tokenMasterEditionPubkey, err := token_metadata.GetMasterEdition(mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid master edition, err: %v", err)
	}

	mintAccountRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), token.MintAccountSize)
	if err != nil {
		log.Fatalf("failed to get mint account rent, err: %v", err)
	}

	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{mint, feePayer},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				system.CreateAccount(system.CreateAccountParam{
					From:     feePayer.PublicKey,
					New:      mint.PublicKey,
					Owner:    common.TokenProgramID,
					Lamports: mintAccountRent,
					Space:    token.MintAccountSize,
				}),
				token.InitializeMint(token.InitializeMintParam{
					Decimals: 0,
					Mint:     mint.PublicKey,
					MintAuth: feePayer.PublicKey,
				}),
				token_metadata.CreateMetadataAccount(token_metadata.CreateMetadataAccountParam{
					Metadata:                tokenMetadataPubkey,
					Mint:                    mint.PublicKey,
					MintAuthority:           feePayer.PublicKey,
					Payer:                   feePayer.PublicKey,
					UpdateAuthority:         feePayer.PublicKey,
					UpdateAuthorityIsSigner: true,
					IsMutable:               true,
					MintData: token_metadata.Data{
						Name:                 "Fake SMS #1355",
						Symbol:               "FSMB",
						Uri:                  "https://34c7ef24f4v2aejh75xhxy5z6ars4xv47gpsdrei6fiowptk2nqq.arweave.net/3wXyF1wvK6ARJ_9ue-O58CMuXrz5nyHEiPFQ6z5q02E",
						SellerFeeBasisPoints: 100,
						Creators: &[]token_metadata.Creator{
							{
								Address:  feePayer.PublicKey,
								Verified: true,
								Share:    100,
							},
						},
					},
				}),
				associated_token_account.Create(associated_token_account.CreateParam{
					Funder:                 feePayer.PublicKey,
					Owner:                  feePayer.PublicKey,
					Mint:                   mint.PublicKey,
					AssociatedTokenAccount: ata,
				}),
				token.MintTo(token.MintToParam{
					Mint:   mint.PublicKey,
					To:     ata,
					Auth:   feePayer.PublicKey,
					Amount: 1,
				}),
				token_metadata.CreateMasterEdition(token_metadata.CreateMasterEditionParam{
					Edition:         tokenMasterEditionPubkey,
					Mint:            mint.PublicKey,
					UpdateAuthority: feePayer.PublicKey,
					MintAuthority:   feePayer.PublicKey,
					Metadata:        tokenMetadataPubkey,
					Payer:           feePayer.PublicKey,
					MaxSupply:       pointer.Get[uint64](0),
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a tx, err: %v", err)
	}

	sig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	fmt.Println(sig)
}
