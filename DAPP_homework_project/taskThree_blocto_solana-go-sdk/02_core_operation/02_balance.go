package main

import (
	"context"
	"fmt"
	"solana-interactor/constant"

	"github.com/blocto/solana-go-sdk/client"
)

func main2() {

	dev_url := constant.DEV_URL
	client := client.NewClient(dev_url)

	//加载私钥
	privateKeyPath := constant.PK_FILE_PATH
	senderAccount, err := loadPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		fmt.Printf("加载私钥失败: %v", err)
	}

	fmt.Printf("发送方地址: %s\n", senderAccount.PublicKey.ToBase58())

	//查询余额（可选）
	balance, err := client.GetBalance(context.Background(), senderAccount.PublicKey.ToBase58())
	if err != nil {
		fmt.Printf("查询发送方余额失败: %v", err)
	} else {
		fmt.Printf("查询发送方余额: %f SOL\n", float64(balance)/1e9)
	}

}
