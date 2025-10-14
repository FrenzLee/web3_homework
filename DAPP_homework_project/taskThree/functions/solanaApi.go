package functions

import (
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
)

/*
 * 将 Lamports（Solana 的最小单位）转换为 SOL（主单位），并保留最多 6 位小数，以字符串形式返回
 * 1 SOL = 10^(9) Lamports
 */
func LamportsToSol(value uint64) string {
	//使用 decimal.New() 直接构造不会损失精度
	sol := decimal.New(int64(value), -9) // value * 10^(-9)
	//截取6位，如果要四舍五入保留6位：sol.Round(6).String()
	return sol.Truncate(6).String()
}

/**
 * 查询账户余额
 */
func GetBalance(rpcClient *rpc.Client, accountBase58 string) (balance string, err error) {
	if rpcClient == nil {
		return balance, errors.New("rpcClient 为空 ！")
	}

	accountPublicKey := solana.MustPublicKeyFromBase58(accountBase58)
	balanceResult, err := rpcClient.GetBalance(context.TODO(), accountPublicKey, rpc.CommitmentFinalized)
	if err != nil {
		return balance, err
	}
	fmt.Println("账户余额：", LamportsToSol(balanceResult.Value))

	return LamportsToSol(balanceResult.Value), nil
}

/*
 * 查询账户信息
 */
func GetAccountInfo(rpcClient *rpc.Client, accountBase58 string) {
	accountPublicKey := solana.MustPublicKeyFromBase58(accountBase58)
	accountInfo, err := rpcClient.GetAccountInfo(context.TODO(), accountPublicKey)
	if err != nil {
		fmt.Println("查询账户信息失败：", err)
	}
	//包括账户数据，Lamports（余额），Owner（账户归属），Data（程序数据）,Executable(是否为可执行程序)，RentEpoch（下一个周期租金-已废弃）,Space（占用空间大小）
	fmt.Println("余额Lamports:", accountInfo.Value.Lamports)
	fmt.Println("余额sol:", LamportsToSol(accountInfo.Value.Lamports))
	fmt.Println("账户归属owner:", accountInfo.Value.Owner.String())
	fmt.Println("占用空间大小space:", accountInfo.Value.Space)
}

/**
 * 获取区块信息
 */
func GetBlockInfo(rpcClient *rpc.Client, blockNum uint64) {
	if blockNum <= 0 {
		//获取最新区块的区块高度
		latestBlockhashResult, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			fmt.Println("获取最新区块信息失败：", err)
		}
		blockNum = latestBlockhashResult.Value.LastValidBlockHeight
	}

	version := uint64(0)
	opts := &rpc.GetBlockOpts{
		Encoding:                       solana.EncodingBase64,
		MaxSupportedTransactionVersion: &version,
	}

	blockInfo, err := rpcClient.GetBlockWithOpts(context.TODO(), blockNum, opts)
	if err != nil {
		fmt.Println("获取区块信息失败：", err)
	}

	//主要包括：blockhash（区块hash），previousBlockhash（上一个区块hash），parentSlot（父存储曹），transactions（交易列表），signatures（交易签名），blockTime（区块确认时间，blockHeight（区块高度））
	fmt.Println("区块信息：", blockInfo)

}
