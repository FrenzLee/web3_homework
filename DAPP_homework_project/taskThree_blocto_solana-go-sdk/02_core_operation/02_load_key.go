package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
)

func main3() {

	//加载私钥
	privateKeyPath := `C:\Users\HP\.config\solana\id.json`
	senderAccount, err := loadPrivateKeyFromFile21(privateKeyPath)
	if err != nil {
		fmt.Printf("加载私钥失败: %v", err)
	}

	fmt.Printf("发送方地址: %s\n", senderAccount.PublicKey.ToBase58())

}

func loadPrivateKeyFromFile21(filepath string) (types.Account, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return types.Account{}, fmt.Errorf("无法读取私钥文件: %v", err)
	}

	var secretKey []byte
	err = json.Unmarshal(data, &secretKey)
	if err != nil {
		return types.Account{}, fmt.Errorf("解析私钥失败: %v", err)
	}

	if len(secretKey) != 64 {
		return types.Account{}, fmt.Errorf("私钥必须是 64 字节，实际: %d", len(secretKey))
	}

	// ✅ 直接使用 blocto SDK 提供的方法
	return types.AccountFromBytes(secretKey)
}

func loadPrivateKeyFromFile22(filepath string) (types.Account, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return types.Account{}, fmt.Errorf("无法读取私钥文件: %v", err)
	}

	var secretKey []byte
	err = json.Unmarshal(data, &secretKey)
	if err != nil {
		return types.Account{}, fmt.Errorf("解析私钥失败: %v", err)
	}

	if len(secretKey) != 64 {
		return types.Account{}, fmt.Errorf("私钥必须是 64 字节，实际: %d", len(secretKey))
	}

	// 前 32 字节是 secret key（seed）
	seed := secretKey[32:]
	// 后 32 字节是公钥
	pubKeyBytes := secretKey[:32]

	return types.Account{
		PrivateKey: seed,
		PublicKey:  common.PublicKeyFromBytes(pubKeyBytes),
	}, nil
}
