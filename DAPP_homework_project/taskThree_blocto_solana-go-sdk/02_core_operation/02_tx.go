package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"solana-interactor/constant"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/types"
)

func main() {

	dev_url := constant.DEV_URL
	client := client.NewClient(dev_url)

	//加载私钥
	privateKeyPath := constant.PK_FILE_PATH
	senderAccount, err := loadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		fmt.Printf("加载私钥失败: %v", err)
	}

	fmt.Printf("发送方地址: %s\n", senderAccount.PublicKey.ToBase58())

	//设置接收方地址
	recipientBase58 := "GCa13vAYRTffkJQNSvwsvudMrGDTUBTuVUNpdMpLSsYb"
	recipientPubkey := common.PublicKeyFromString(recipientBase58)

	//查询余额（可选）
	balance, err := client.GetBalance(context.Background(), senderAccount.PublicKey.ToBase58())
	if err != nil {
		fmt.Printf("查询发送方余额失败: %v", err)
	} else {
		fmt.Printf("查询发送方余额: %f SOL\n", float64(balance)/1e9)
	}

	//构建交易
	latestBlockHash, err := client.GetLatestBlockhash(context.Background())
	if err != nil {
		fmt.Printf("获取latestBlockHash失败: %v", err)
	}

	//创建转账指令（转账0.01 SOL）
	instruction := system.Transfer(system.TransferParam{
		From:   senderAccount.PublicKey,
		To:     recipientPubkey,
		Amount: 100_000_000, // 0.1 SOL
	})

	//构建Message
	message := types.NewMessage(types.NewMessageParam{
		FeePayer:        senderAccount.PublicKey,
		RecentBlockhash: latestBlockHash.Blockhash,
		Instructions:    []types.Instruction{instruction},
	})

	//创建交易
	transaction, err := types.NewTransaction(
		types.NewTransactionParam{
			Message: message,
			Signers: []types.Account{senderAccount},
		},
	)
	if err != nil {
		fmt.Printf("创建交易失败: %v", err)
	}

	//发送交易
	txHash, err := client.SendTransaction(context.Background(), transaction)
	if err != nil {
		fmt.Printf("发送交易失败: %v", err)
	}

	//输出结果
	fmt.Printf("交易哈希: %s\n", txHash)

}

// 从文件加载私钥（推荐方式）
func loadPrivateKeyFromFile(filepath string) (types.Account, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return types.Account{}, fmt.Errorf("无法读取私钥文件: %v", err)
	}

	var secretKey []byte
	err = json.Unmarshal(data, &secretKey)
	if err != nil {
		return types.Account{}, fmt.Errorf("解析私钥失败: %v", err)
	}

	if len(secretKey) != 64 {
		return types.Account{}, fmt.Errorf("私钥必须是 64 字节，实际: %d", len(secretKey))
	}

	// 直接使用 blocto SDK 提供的方法
	return types.AccountFromBytes(secretKey)
}
