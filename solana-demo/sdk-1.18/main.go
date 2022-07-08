package main

import (
	"context"
	"encoding/json"
	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
	"log"
)

func main() {
	//account()
	//transfer(1)
	//erc20Balance()
	erc20Transfer()
	//getAirdrop()
}

func transfer(amount uint64) (hash string, err error) {
	c := client.NewClient(rpc.DevnetRPCEndpoint)
	wallet, err := types.AccountFromBase58("41zA8sti5SjxbzN4e23LTUUcoTaNfrbM6jAQWYvBhpgV7U2M28PYxQ96KV2FifM2ZqpscxSaQgeGHG1NtQhLewzg")
	if err != nil {
		log.Println("AccountFromBase58 err: ", err)
		return "", err
	}
	log.Println("wallet address: ", wallet.PublicKey.String())

	to := common.PublicKeyFromString("4jHT5F3DaUr2ZVAtQsTngNZ1LroTb1B3dY6ARjBd32jL")
	log.Println("to address: ", "4jHT5F3DaUr2ZVAtQsTngNZ1LroTb1B3dY6ARjBd32jL")

	res, err := c.GetLatestBlockhash(context.TODO())
	if err != nil {
		log.Println("GetRecentBlockhash err: ", err)
		return "", err
	}

	amount = amount * 1e6
	msgParam := types.NewMessageParam{
		FeePayer:        wallet.PublicKey,
		Instructions:    []types.Instruction{
			sysprog.Transfer(sysprog.TransferParam{
				From:   wallet.PublicKey,
				To:     to,
				Amount: amount,
			}),
		},
		RecentBlockhash: res.Blockhash,
	}
	message := types.NewMessage(msgParam)
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: []types.Account{
			{
				PublicKey:  wallet.PublicKey,
				PrivateKey: wallet.PrivateKey,
			},
		},
	})

	txHash, err := c.SendTransaction(context.TODO(), tx)
	if err != nil {
		log.Println("SendTransaction err: ", err)
		return "", err
	}

	log.Println("hash: ", txHash)
	return txHash, nil
}

func erc20Balance() {

	// SPL-TOKEN账户可直接查余额
	account := "GDSNZP4oUygNNitjyUJr4Y4TJKouTZdfQDopfSgUEytu"
	rc := rpc.NewRpcClient(rpc.DevnetRPCEndpoint)
	res, err := rc.GetTokenAccountBalance(context.Background(), account)
	if err != nil {
		log.Fatal("GetTokenAccountBalance err", err)
	}

	b, _ := json.Marshal(res.Result)
	log.Println(string(b))

	// 查询SOL余额
	rc.GetBalance(context.TODO(), account)

	// SPL-TOKEN账户可直接查余额，返回余额和精度
	c := client.NewClient(rpc.DevnetRPCEndpoint)
	balance, decaimal, err := c.GetTokenAccountBalance(context.Background(), account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Balance: ", balance)
	log.Println("Decimal: ", decaimal)

	// 查询账户代币余额
	account2 := "DW3dfsVrjEZqJXyfq7uzCMfwUD35rKeWzX6H8tj22fpk"
	programId := "7WCgirvbCkoTjQrKEAZStcmJYfbZ5qPuYnRxLgshnooD"
	accountsRes, err := rc.GetTokenAccountsByOwnerWithConfig(
		context.TODO(),
		account2,
		rpc.GetTokenAccountsByOwnerConfigFilter{
			Mint: programId,
		},
		rpc.GetTokenAccountsByOwnerConfig{Encoding: "jsonParsed"},
	)
	if err != nil {
		log.Fatal(err)
	}
	rb, _ := json.Marshal(accountsRes.Result)
	log.Println(string(rb))
}

func erc20Transfer() (hash string, err error) {
	c := client.NewClient(rpc.DevnetRPCEndpoint)
	wallet, err := types.AccountFromBase58("41zA8sti5SjxbzN4e23LTUUcoTaNfrbM6jAQWYvBhpgV7U2M28PYxQ96KV2FifM2ZqpscxSaQgeGHG1NtQhLewzg")
	if err != nil {
		log.Fatal("AccountFromBase58 err: ", err)
		return "", err
	}
	log.Println("wallet address: ", wallet.PublicKey.String())

	to := common.PublicKeyFromString("BDFLaeF72qKry2f35grXfxBGQJRWnGzSZDdhYddvHVJs")
	log.Println("to address: ", to.ToBase58())

	token := common.PublicKeyFromString("7WCgirvbCkoTjQrKEAZStcmJYfbZ5qPuYnRxLgshnooD")
	log.Println("token address", token.ToBase58())

	res, err := c.GetLatestBlockhash(context.TODO())
	if err != nil {
		log.Fatal("GetRecentBlockhash err: ", err)
		return "", err
	}
	log.Println("BlockHash:", res.Blockhash)

	fromATA, _, err := common.FindAssociatedTokenAddress(wallet.PublicKey, token)
	if err != nil {
		log.Fatal("FindAssociatedTokenAddress wallet err: ", err)
	}
	log.Println("FindAssociatedTokenAddress", "fromATA", fromATA.ToBase58())

	toATA, _, err := common.FindAssociatedTokenAddress(to, token)
	if err != nil {
		log.Fatal("FindAssociatedTokenAddress to err: ", err)
	}
	log.Println("FindAssociatedTokenAddress", "toATA", toATA.ToBase58())

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        wallet.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions:    []types.Instruction{
				tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
					From:     fromATA,
					To:       toATA,
					Mint:     token,
					Auth:     wallet.PublicKey,
					Signers:  []common.PublicKey{},
					Amount:   1e8,
					Decimals: 9,
				}),
			},
		}),
		Signers: []types.Account{wallet},
	})

	txHash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal("SendTransaction err: ", err)
		return "", err
	}
	log.Println("txHash: ", txHash)
	return "", nil
}

func getAirdrop() {
	rc := rpc.NewRpcClient(rpc.DevnetRPCEndpoint)
	res, err := rc.RequestAirdrop(context.TODO(), "DW3dfsVrjEZqJXyfq7uzCMfwUD35rKeWzX6H8tj22fpk", 2 * 1e9)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
}

func account() {
	var prvBytes = []byte{32,205,107,113,25,29,7,50,188,114,185,227,26,134,1,123,77,135,46,249,241,212,84,201,79,142,34,6,29,204,253,35,151,184,12,239,196,53,60,137,252,102,19,57,51,108,166,113,133,153,126,182,241,126,251,26,159,47,209,220,225,113,235,204}
	account, err := types.AccountFromBytes(prvBytes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(account.PublicKey.ToBase58())
	log.Println(base58.Encode(prvBytes))

}