package functions

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func QueryBlockInfo(client *ethclient.Client, blockNumber int64) {

	block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		log.Fatal("get block err: ", err)
	}

	fmt.Printf("区块号: %d\n", block.Number().Uint64())
	fmt.Printf("区块时间戳: %d\n", block.Time())
	fmt.Printf("区块难度: %d\n", block.Difficulty().Uint64())
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("区块交易数量: %d\n", len(block.Transactions()))

}
