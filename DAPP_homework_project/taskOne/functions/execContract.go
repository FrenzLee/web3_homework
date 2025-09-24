package functions

import (
	"context"
	"fmt"
	"log"

	"github.com/dapp_homework/taskOne/constant"
	"github.com/dapp_homework/taskOne/functions/counter"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ExecContract(client *ethclient.Client) {

	const contractAddrHexStr = constant.ContractAddr
	contractAddress := common.HexToAddress(contractAddrHexStr)

	//获取合约实例
	counterContract, err := counter.NewCounter(contractAddress, client)
	if err != nil {
		log.Fatal("get counterContract err: ", err)
	}

	//加载私钥
	privateKeyHexStr := constant.PK
	privatekey, err := crypto.HexToECDSA(privateKeyHexStr)
	if err != nil {
		log.Fatal("私钥加载失败：", err)
	}

	//获取chainID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("get chainID err ：", err)
	}

	//创建交易签名器
	auth, err := bind.NewKeyedTransactorWithChainID(privatekey, chainID)
	if err != nil {
		log.Fatal("get auth err ：", err)
	}

	//调用合约方法
	tx, err := counterContract.IncrementOne(auth)
	if err != nil {
		log.Fatal("tx err ：", err)
	}
	fmt.Printf("交易哈希：%s \n", tx.Hash().Hex())

	//获取当前count
	result, err := counterContract.GetCurrentCount(&bind.CallOpts{})
	if err != nil {
		log.Fatal("get result err ：", err)
	}
	fmt.Printf("当前count：%s \n", result)

}
