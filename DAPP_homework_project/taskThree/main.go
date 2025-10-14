package main

import (
	"Dapp-solana-taskThree/constant"
	"Dapp-solana-taskThree/functions"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

func main() {

	//cluster := rpc.DevNet
	rpcClient := rpc.NewWithCustomRPCClient(rpc.NewWithLimiter(constant.DEV_URL, rate.Every(time.Second), 5))
	defer rpcClient.Close()

	account := "HW5Sjq3wPkUMJwjzuLhBoQnXszTDyt739ccYrGN4NgpC"

	//查询账户余额
	_, err := functions.GetBalance(rpcClient, account)
	if err != nil {
		fmt.Printf("❌ 查询余额失败: %v\n", err)
		return
	}

	//查询账户信息
	//functions.GetAccountInfo(rpcClient, account)

	//获取区块信息
	//functions.GetBlockInfo(rpcClient, 0)

	/*
		//交易转账
		payerPrivateKey := constant.FROM_ACCOUNT_PRIVATE_KEY
		toAccount := constant.TO_ACCOUNT

		wsClient, err := ws.Connect(context.TODO(), constant.DEV_URL_WS)
		if err != nil {
			fmt.Println("Failed to connect to ws: ", err)
		}
		defer wsClient.Close()

		//轮询RPC获取交易结果
		//functions.SendTransfer(rpcClient, payerPrivateKey, toAccount, 1)

		//使用 WebSocket获取交易结果
		//functions.SendTransferAndSubscript(rpcClient, wsClient, payerPrivateKey, toAccount, 1)

		fmt.Printf("转账后，查询账户余额")
		functions.GetBalance(rpcClient, account)
	*/

	//查询交易详情
	signature := "j5p9w83tXN7CuCQeoiXg7j7WRBT6suVd9ZyHYr54PHHm2rGTeYGxHBqdgr4WLX8WwzFx5SKG6wzN2MME47XZWYD"
	functions.GetTransactionInfo(rpcClient, signature)

}
