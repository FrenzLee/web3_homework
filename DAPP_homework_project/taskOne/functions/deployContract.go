package functions

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/dapp_homework/taskOne/constant"
	"github.com/dapp_homework/taskOne/functions/counter"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func DeployContract(client *ethclient.Client) {

	//加载私钥
	privateKeyHexStr := constant.PK
	privatekey, err := crypto.HexToECDSA(privateKeyHexStr)
	if err != nil {
		log.Fatal("私钥加载失败：", err)
	}

	//获取公钥，推导出转账付款账户地址
	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("get publicKeyECDSA err：", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("合约部署账户地址：%s \n", fromAddress)

	//获取nonce（交易计数器）
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("get nonce err：", err)
	}

	//获取 Gas Price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("get gasPrice err ：", err)
	}

	//获取chainID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("get chainID err ：", err)
	}

	//创建交易签名器
	auth, err := bind.NewKeyedTransactorWithChainID(privatekey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	//部署智能合约
	address, tx, instance, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatal("store DeployStore err ：", err)
	}
	_ = instance

	fmt.Printf("交易哈希：%s \n", tx.Hash().Hex())
	//0x42b34aca25c0fe5f528227a23d832da47f9055fe0a8771bac3853ae695d21847
	
	fmt.Printf("合约地址：%s \n", address.Hex())
	//0x94EEbce00e7EA840f17c9FCca13675633C3164f2

}
