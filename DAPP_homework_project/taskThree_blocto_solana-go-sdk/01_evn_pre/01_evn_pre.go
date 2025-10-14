package main

import (
	"context"
	"fmt"
	"solana-interactor/constant"

	"github.com/blocto/solana-go-sdk/client"
)

func main() {

	c := client.NewClient(constant.DEV_URL)

	version, err := c.GetVersion(context.Background())
	if err != nil {
		fmt.Println("❌ GetVersion failed:", err)
		return
	}

	println("Solana node version:", version.SolanaCore)

}
