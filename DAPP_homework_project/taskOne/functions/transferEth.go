package functions

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/dapp_homework/taskOne/constant"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TransferEth(client *ethclient.Client, toAddressHexStr string) {
	//加载私钥
	privateKeyHexStr := constant.PK
	privateKey, err := crypto.HexToECDSA(privateKeyHexStr)
	if err != nil {
		log.Fatal("私钥加载失败：", err)
	}

	//获取公钥，推导出转账付款账户地址
	publicKey := privateKey.Public()                   //返回值是接口类型
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey) //类型断言转换，value, ok := interface.(ConcreteType)
	if !ok {
		log.Fatal("获取公钥失败：", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("获取到fromAddress：", fromAddress.Hex())

	//获取nonce（交易计数器）
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("获取nonce失败：", err)
	}

	//获取预估gasPrice
	tip, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal("获取tip失败：", err)
	}
	//获取最新区块baseFee
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal("获取header失败：", err)
	}
	baseFee := header.BaseFee
	//计算GasFeeCap
	gasFeeCap := new(big.Int).Add(baseFee, tip)

	//设置交易参数
	value := big.NewInt(100000000000000000) //18个0是 1 eth, 17个0是 0.1 eth。目前是 0.1 eth
	gasLimit := uint64(21000)               //标准转账 Gas 限制
	toAddress := common.HexToAddress(toAddressHexStr)
	var data []byte //附加数据，转账时为空

	//获取chainID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("获取chainID失败：", err)
	}

	//构造交易
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &toAddress,
		Value:     value,
		Gas:       gasLimit,
		GasTipCap: tip,
		GasFeeCap: gasFeeCap,
		Data:      data,
	})

	//签名交易
	signer := types.NewLondonSigner(chainID)
	signerTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatal("签名交易失败：", err)
	}

	//广播交易
	err = client.SendTransaction(context.Background(), signerTx)
	if err != nil {
		log.Fatal("广播交易失败：", err)
	}

	//输出交易哈希
	fmt.Printf("交易哈希 tx sent: %s", signerTx.Hash().Hex())
	//0x24eb0621985f288111c8453fe246f7bef42be2b32da3168568e44372c3c33f0a

}
