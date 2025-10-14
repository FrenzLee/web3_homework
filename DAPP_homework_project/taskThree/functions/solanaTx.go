package functions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	sendandconfirmtransaction "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

/*
 * 交易
 */
func SendTransfer(rpcClient *rpc.Client, payerPrivateKey string, toAccount string, amount uint64) (*solana.Signature, error) {
	privateKey, err := solana.PrivateKeyFromBase58(payerPrivateKey)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.PublicKey()
	toAccountPublicKey := solana.MustPublicKeyFromBase58(toAccount)

	lamports := amount * 1e9
	if lamports <= 0 {
		return nil, errors.New("转账金额必须大于0")
	}

	recentBlockHash, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	//构建交易
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				lamports,           //金额
				publicKey,          //转出账户
				toAccountPublicKey, //转入账户
			).Build(),
		},
		recentBlockHash.Value.Blockhash,
		solana.TransactionPayer(publicKey),
	)
	if err != nil {
		fmt.Println("构建交易失败：", err)
		return nil, err
	}

	//签名
	tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key == publicKey {
			return &privateKey
		}
		return nil
	})

	//发送交易
	signature, err := rpcClient.SendTransaction(context.TODO(), tx)
	if err != nil {
		fmt.Println("发送交易失败：", err)
		return nil, err
	}

	fmt.Println("发送交易成功,签名：", signature)

	for i := 0; i < 20; i++ {
		res := GetTxStatus(rpcClient, signature)
		if res == 1 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return &signature, nil
}

/*
 * 查询交易结果
 */
func GetTxStatus(rpcClient *rpc.Client, signature solana.Signature) int64 {
	result, _ := rpcClient.GetSignatureStatuses(context.TODO(), true, signature)

	if result.Value != nil && len(result.Value) > 0 {
		switch result.Value[0].ConfirmationStatus {
		case rpc.ConfirmationStatusFinalized:
			fmt.Println("交易已最终确认")
			return 1
		case rpc.ConfirmationStatusConfirmed:
			return 2
		case rpc.ConfirmationStatusProcessed:
			return 0
		default:
			return 0
		}
	}

	return -1
}

/*
 * 发起交易 通过websocket监听状态，流程与常规交易一致
 */
func SendTransferAndSubscript(rpcClient *rpc.Client, wsClient *ws.Client, payerPrivateKey string, toAccount string, amount uint64) (*solana.Signature, error) {
	privateKey := solana.MustPrivateKeyFromBase58(payerPrivateKey)
	fromAccountPublicKey := privateKey.PublicKey()
	toPublickey := solana.MustPublicKeyFromBase58(toAccount)

	lamports := amount * 1e9
	if lamports <= 0 {
		return nil, errors.New("转账金额必须大于0")
	}

	recentBlockHash, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	//构建交易
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				lamports,
				fromAccountPublicKey,
				toPublickey,
			).Build(),
		},
		recentBlockHash.Value.Blockhash,
		solana.TransactionPayer(fromAccountPublicKey),
	)
	if err != nil {
		fmt.Println("构建交易失败：", err)
		return nil, err
	}

	//签名
	tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(fromAccountPublicKey) {
			return &privateKey
		}
		return nil
	})

	//发送交易
	signature, err := sendandconfirmtransaction.SendAndConfirmTransaction(context.TODO(), rpcClient, wsClient, tx)
	if err != nil {
		fmt.Println("发送交易失败：", err)
		return nil, err
	}
	fmt.Println("发送交易成功,签名：", signature)

	//订阅交易状态
	subscribe, err := wsClient.SignatureSubscribe(signature, rpc.CommitmentFinalized)
	defer subscribe.Unsubscribe()

	for {
		select {
		case <-subscribe.Response():
			fmt.Println("交易已确认")
			break
		case <-time.After(time.Second * 30):
			fmt.Println("交易确认超时")
			break
		}
	}

	return &signature, nil
}

/*
 * 查询交易详情
 */
func GetTransactionInfo(rpcClient *rpc.Client, signature string) {
	version := uint64(0)
	TxInfo, err := rpcClient.GetTransaction(
		context.TODO(),
		solana.MustSignatureFromBase58(signature),
		&rpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			MaxSupportedTransactionVersion: &version,
		},
	)
	if err != nil {
		fmt.Println("获取交易详情出错：", err)
		return
	}

	fmt.Println("区块高度：", TxInfo.Slot)
	fmt.Println("区块时间：", TxInfo.BlockTime.String())
	TransactionMeta := TxInfo.Meta

	fmt.Println("手续费：", TransactionMeta.Fee)
	fmt.Println("交易前余额：", TransactionMeta.PreBalances[0])
	fmt.Println("交易后余额：", TransactionMeta.PostBalances[0])
	fmt.Println("交易日志：", TransactionMeta.LogMessages)
	TransactionResultEnvelope := TxInfo.Transaction

	Transaction, err := TransactionResultEnvelope.GetTransaction()
	fmt.Println("最近一个区块hash：", Transaction.Message.RecentBlockhash.String())
	MessageHeader := Transaction.Message.Header
	fmt.Println("MessageHeader:", MessageHeader)
	fmt.Println("Transaction->IsVersioned:", Transaction.Message.IsVersioned())
	fmt.Println("交易账户：", Transaction.Message.AccountKeys)
	tblKeys := Transaction.Message.GetAddressTableLookups().GetTableIDs()
	fmt.Println("Transaction->tblKeys:", tblKeys)
	for i, instruction := range Transaction.Message.Instructions {
		fmt.Printf("指令 %d: 程序ID索引=%d, 数据长度=%d\n",
			i, instruction.ProgramIDIndex, len(instruction.Data))
	}

}
