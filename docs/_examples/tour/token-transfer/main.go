package main

import (
	"context"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

// DW3dfsVrjEZqJXyfq7uzCMfwUD35rKeWzX6H8tj22fpk
var feePayer, _ = types.AccountFromBase58("41zA8sti5SjxbzN4e23LTUUcoTaNfrbM6jAQWYvBhpgV7U2M28PYxQ96KV2FifM2ZqpscxSaQgeGHG1NtQhLewzg")

var mintPubkey = common.PublicKeyFromString("7WCgirvbCkoTjQrKEAZStcmJYfbZ5qPuYnRxLgshnooD")

var fromATA = common.PublicKeyFromString("A5g7ja5fHA1jnXYwhUUGRoVxEj4MA4P76Bw6ggu4rUu9")

var toATA = common.PublicKeyFromString("GDSNZP4oUygNNitjyUJr4Y4TJKouTZdfQDopfSgUEytu")

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	res, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
					From:     fromATA,
					To:       toATA,
					Mint:     mintPubkey,
					Auth:     feePayer.PublicKey,
					Signers:  []common.PublicKey{},
					Amount:   1e8,
					Decimals: 9,
				}),
			},
		}),
		Signers: []types.Account{feePayer},
	})
	if err != nil {
		log.Fatalf("failed to new tx, err: %v", err)
	}

	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send raw tx error, err: %v\n", err)
	}

	log.Println("txhash:", txhash)
}