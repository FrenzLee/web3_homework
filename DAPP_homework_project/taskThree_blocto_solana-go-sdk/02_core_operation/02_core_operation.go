package main

import (
	"context"
	"fmt"
	"solana-interactor/constant"

	"github.com/blocto/solana-go-sdk/client"
)

func main1() {
	dev_url := constant.DEV_URL
	client := client.NewClient(dev_url)

	//获取最新区块的slot
	slot, err := client.GetSlot(context.Background())
	if err != nil {
		fmt.Println("获取最新区块slot失败：", err)
		return
	}
	fmt.Printf("最新slot: %d\n", slot)

	// 获取最新区块
	block, err := client.GetBlock(context.Background(), uint64(slot))
	if err != nil {
		fmt.Println("获取最新区块信息失败：", err)
		return
	}

	fmt.Printf("区块高度: %d\n", block.BlockHeight)
	fmt.Printf("交易数量: %d\n", len(block.Transactions))
	fmt.Printf("区块Hash: %s\n", block.Blockhash)

}
