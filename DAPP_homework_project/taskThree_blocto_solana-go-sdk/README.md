### 一. Solana-go 相关资料
```
官方文档：https://pkg.go.dev/github.com/gagliardetto/solana-go#section-readme
源码地址：https://github.com/gagliardetto/solana-go?tab=readme-ov-file#features
浏览器：https://solscan.io/?cluster=devnet
测试水领取地址：https://faucet.solana.com/
钱包：Phantom、Backpack、Solflare
```

### 二. Solana交易流程
```
客户端发送PRC请求 
 ---> 创建交易 ---> 构造交易数据 ---> 签名 ---> 交易广播 
 ---> 节点预验证（数据正确性） ---> 领导节点打包交易生成区块 
 ---> 其他共识节点验证 ---> 上链 ---> 交易确认 
```

###  三. 项目使用SDK
github.com/blocto/solana-go-sdk

###  四. github.com/blocto/solana-go-sdk 和 github.com/gagliardetto/solana-go
github.com/blocto/solana-go-sdk 和 github.com/gagliardetto/solana-go 是两个流行的 Go 语言 SDK，用于与 Solana 区块链进行交互。虽然它们目标相似，但在设计理念、功能覆盖、维护方向和使用场景上有显著区别。
| 特性/项目 | `blocto/solana-go-sdk` | `gagliardetto/solana-go` |
|----------|------------------------|---------------------------|
| **维护者** | Blocto（钱包团队，支持 Phantom、Sollet 等） | Alessandro Gagliardetto（独立开发者，社区知名） |
| **定位** | 轻量级、面向应用开发（如钱包、DApp 后端） | 全功能、面向开发者和工具链（CLI、解析器等） |
| **API 设计** | 更简洁、易用，偏向高层封装 | 更底层、灵活，贴近 Solana 官方 RPC 和 BPF 语义 |
| **BPF 程序支持** | 提供常见代币（SPL Token）、Stake、System 等指令封装 | 提供完整的指令生成 + 解析（可反序列化链上程序数据） |
| **RPC 客户端** | 封装了常用 RPC 调用，API 友好 | 功能最全，支持几乎所有 Solana RPC 方法 |
| **交易构建** | 简化交易构建流程，适合钱包集成 | 完整支持 V0、Legacy 交易，支持地址查找表（Address Lookup Tables） |
| **序列化/反序列化** | 有限支持（主要关注发送交易） | 强大支持：可解析链上账户数据（如代币账户、Stake 账户等） |
| **文档** | 一般，依赖代码示例 | 较好，有 godoc 和 CLI 工具示例 |
| **生态集成** | 常用于钱包、DApp、Blocto 生态 | 常用于分析工具、索引器、浏览器、CLI 工具 |
| **是否支持 WASM** | ✅ 支持编译为 WASM（前端可用） | ❌ 不支持 WASM |
| **社区活跃度** | 中等 | 非常高（GitHub Star 超 1.5k，广泛引用） |
| **License** | MIT | MIT |

如何选择？

| 你的需求 | 推荐 SDK |
|--------|---------|
| 快速开发 DApp、钱包、后端服务 | ✅ `blocto/solana-go-sdk` |
| 需要在浏览器或前端使用 Go 编译的代码 | ✅ `blocto/solana-go-sdk`（支持 WASM） |
| 做链上数据分析、解析账户、构建索引器 | ✅ `gagliardetto/solana-go` |
| 需要解析 Raydium、Orca、Marinade 等协议数据 | ✅ `gagliardetto/solana-go` |
| 构建 CLI 工具或自动化脚本 | ✅ `gagliardetto/solana-go` |
| 与 Blocto、Phantom 钱包集成 | ✅ `blocto/solana-go-sdk` |



