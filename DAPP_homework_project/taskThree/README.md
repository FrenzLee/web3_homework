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

###  三. Solana账户模型
区别于比特币的UTXO和以太坊的EVM的账户模型，Solana账户程序与存储的数据分离。
#### 1. 基本概念
区别于以太坊EVM的EOA和合约账户，Solana中每个账户由以下几个核心属性组成：<br>
（1）余额  lamports <br>
（2）字节数据 data <br>
（3）账户归属应用 owner <br>
（4）是否为可执行程序 is_executable <br>
（5）下一个周期租金  rent_epoch <br>
账户可存放可执行的代码程序、程序运行的状态数据，在solana中数据存储需要有一定的租金，否则数据将被清理。
#### 2. 账户归属与交易
（1）账户的归属权限中，只有owner（指定的programId）才能写入、修改该账户的data数据，也就是说修改数据必须指定owner，读取则不需要。<br>
（2）发起交易必须声明所有涉及访问的数据，由instruction构成。instruction中指明哪些账户读写，哪些要签名等，这与evm区别很大，也正是因为这个机制，solana能支持大量并发执行。<br>
（3）数据存储租金，账户数据需要支付一定金额租金才能将数据保存在链上，一般可以通过RPC的【GetMinimumBalanceForRentExemption】获取数据所需要支付的最低金额。<br>
#### 3. 交易并行
（1）solana在交易提交时需要明确申明涉及的读写账户，SVM会把不冲突的交易打包并行执行，从而实现高吞吐。其中修改账户启用排它锁，只读可并行执行，以提高TPS。
#### 4. 与EVM区别
| 维度        |                                           EVM(ETH)                                            |                        SVM(Solana)                  
                  |
|:----------|:---------------------------------------------------------------------------------------------:|------------------------------------------------------------------------------------------------------------------------:|
| 账户类型      | EVM账户主要分2种：EOA钱包账户和合约账户。<br/>  • EOA钱包账户由私钥控制，可用发起交易。<br/> • 合约账户由EOA发起部署，所有人为发起账户，合约程序与存储耦合。 | 账户程序分离模型。包括程序账户与数据账户。<br/> • 数据账户由系统程序创建，存储程序运行的状态与数据。<br/> • 程序账户由BPF加载器创建。新的数据账户创建后会将所有权转移给该程序，程序账户标记executable=true。 |
| 账户权限      |                                EOA由私钥控制，合约由代码控制（内部实现owner逻辑）。                                 |                                                                                   数据账户的owner归属程序program，修改数据账户只能程序签名修改。 |
| 状态存储      |                                         合约程序与数据存储在一块。                                         |                                                                                                              数据与程序分开存储。 |
| 交易模型      |                              输入from、to、balance，EVM运行时按顺序动态解析决定。                               |                                                                                            交易必须完整申明账户列表，根据账号冲突分组并行打包交易。 |
| 费用模型      |              以gas为单位计费，一笔交易计算与存储一次性计费（基础费21000+小费）<br/> • 基础费用一部分销毁，一部分给矿工<br/> • 小费全部给矿工              |                                         以CU（Compute Unit）计算单元计费，包括交易计算和数据存储费用。<br/> • 计算费用主要由固定基础费（5000）+可选的优先费用。<br/> • 数据存储租金定期付费，可通过最低押金金方式获取数据租金豁免。 |
| 可升级性      | 本身不支持，需要通过call和delegatecall代理实现，将逻辑与数据存储分离。 | 支持通过Upgradeable BPF Loader升级标记为非final的程序。 |
| 并发模型      | 单线程顺序执行。 | 根据申明账户冲突判断，进行分组并行执行，账户锁。 |
#### 5.SVM交易模型账号申明
首先这里所说的交易需要申明所有账户以及访问状态，是因为账户模型结构（程序与数据分离）要求必须如此。
在ETH中一笔交易，合约中往往只需要传入 from,to,amount 即可，剩下的（被调用合约，中间跨合约再调用，代币接受者）都是在执行过程中动态解析的。 
而SVM需要事先声明，程序需要知道它读写哪些账号，需要访问哪些数据单元，否则无法并行执行。 
#### 6.关于SPL Token代币模型说明
在ETH中（ERC20标准），每个token都是一个独立的合约，用户token数量直接在合约中映射，即用户钱包地址==>token数量，简单明了。
Solana的SPL token模型，程序（合约）不存储数据，每个账户自己持有余额数据。其中涉及几个概念，包括token program，mint account，token address，owner
（1）Token program：token 程序（类似ERC20的合约），管理代币逻辑，如发行和转账等。
（2）Mint Account：某个代币类型（元数据），存储totalSupply，decimals ，mint_authority(谁能增发)，类似ERC20中的全局变量。
（3）Token Account：保存某个用户对mint的余额，类似ERC20的balances(address)。
（4）Owner：这个token account 属于谁（钱包地址）。
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



