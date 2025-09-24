package main

import (
	"log"

	"github.com/dapp_homework/taskOne/functions"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/d91b3d846cd64786b470071dd9918416")
	if err != nil {
		log.Fatal("连接eth客户端失败：", err)
	}
	defer client.Close()

	//查询区块信息
	//functions.QueryBlockInfo(client, 9239148)

	//转账
	//functions.TransferEth(client, "0x0c46c27D5642516B833e617566F42Cc302fA33f1")

	//部署Counter合约
	//functions.DeployContract(client)

	//执行合约
	functions.ExecContract(client)

}
