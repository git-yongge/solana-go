package main

import (
	"context"
	"encoding/base64"
	"github.com/gagliardetto/solana-go/programs/token"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {

	endpoint := rpc.DevNet_RPC
	client := rpc.New(endpoint)

	wallet, _ := solana.PrivateKeyFromBase58("41zA8sti5SjxbzN4e23LTUUcoTaNfrbM6jAQWYvBhpgV7U2M28PYxQ96KV2FifM2ZqpscxSaQgeGHG1NtQhLewzg")
	to, _ := solana.PublicKeyFromBase58("BDFLaeF72qKry2f35grXfxBGQJRWnGzSZDdhYddvHVJs")
	mint, _ := solana.PublicKeyFromBase58("7WCgirvbCkoTjQrKEAZStcmJYfbZ5qPuYnRxLgshnooD")
	erc20Transfer(client, wallet.PublicKey(), to, mint, wallet)
}

func erc20Transfer(client *rpc.Client, from, to, mint solana.PublicKey, prv solana.PrivateKey) {

	fromATA, _, _ := solana.FindAssociatedTokenAddress(from, mint)
	toATA, _, _ := solana.FindAssociatedTokenAddress(to, mint)
	log.Println("from:", from)
	log.Println("toATA:", toATA)
	log.Println("mint:", mint)

	transferCheck := token.NewTransferCheckedInstructionBuilder()
	transferCheck.SetSourceAccount(fromATA)
	transferCheck.SetDestinationAccount(toATA)
	transferCheck.SetMintAccount(mint)
	transferCheck.SetOwnerAccount(from)
	transferCheck.SetDecimals(9)
	transferCheck.SetAmount(1e8)

	//instruction := []solana.Instruction{
	//	transferCheck.Build(),
	//}

	res, err := client.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("blockHash:", res.Value.Blockhash)

	builder := solana.NewTransactionBuilder()
	builder.SetRecentBlockHash(res.Value.Blockhash)
	builder.SetFeePayer(from)
	builder.AddInstruction(transferCheck.Build())

	tx, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	//tx, err := solana.NewTransaction(instruction, res.Value.Blockhash)
	//if err != nil {
	//	log.Fatal(err)
	//}
	tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		return &prv
	})
	err = tx.VerifySignatures()
	log.Println("VerifySignatures:", err)

	sign, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("sign:", sign.String())

	txRaw, err := base64.StdEncoding.DecodeString(sign.String())
	if err != nil {
		log.Fatal("DecodeString err:", err)
	}

	sig, err := client.SendRawTransaction(context.TODO(), txRaw)
	if err != nil {
		log.Fatal("SendRawTransaction:", err)
	}
	log.Println(sig.String())
}

func getKey() (solana.PrivateKey, solana.PublicKey, error) {
	base58Key := "41zA8sti5SjxbzN4e23LTUUcoTaNfrbM6jAQWYvBhpgV7U2M28PYxQ96KV2FifM2ZqpscxSaQgeGHG1NtQhLewzg"
	prv, err := solana.PrivateKeyFromBase58(base58Key)
	if err != nil {
		return nil, solana.PublicKey{}, err
	}
	return prv, prv.PublicKey(), nil
}