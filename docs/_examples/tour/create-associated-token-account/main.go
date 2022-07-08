package main

import (
	"context"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/assotokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

// FUarP2p5EnxD66vVDL4PWRoWMzA56ZVHG24hpEDFShEz
var feePayer, _ = types.AccountFromBase58("f3CR8UdWZ6PgXpjwZh6Ak2b3irGFJ9giXMavGeC2vJGdFjEcrmJahCQHj45jKWSFkdjJSSTESqKkjAnVQgd267u")

// 9aE476sH92Vz7DMPyq5WLPkrKWivxeuTKEFKd2sZZcde
var alice, _ = types.AccountFromBase58("41zA8sti5SjxbzN4e23LTUUcoTaNfrbM6jAQWYvBhpgV7U2M28PYxQ96KV2FifM2ZqpscxSaQgeGHG1NtQhLewzg")

var mintPubkey = common.PublicKeyFromString("7WCgirvbCkoTjQrKEAZStcmJYfbZ5qPuYnRxLgshnooD")

func main() {
	c := client.NewClient(rpc.DevnetRPCEndpoint)

	ata, _, err := common.FindAssociatedTokenAddress(alice.PublicKey, mintPubkey)
	if err != nil {
		log.Fatalf("find ata error, err: %v", err)
	}
	fmt.Println("ata:", ata.ToBase58())

	res, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{
					Funder:                 feePayer.PublicKey,
					Owner:                  alice.PublicKey,
					Mint:                   mintPubkey,
					AssociatedTokenAccount: ata,
				}),
			},
		}),
		Signers: []types.Account{feePayer},
	})
	if err != nil {
		log.Fatalf("generate tx error, err: %v\n", err)
	}

	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send raw tx error, err: %v\n", err)
	}

	log.Println("txhash:", txhash)
}