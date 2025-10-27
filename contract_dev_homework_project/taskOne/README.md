## 一、项目概述
`MemeToken` 智能合约，涵盖代币创建、交易、添加/移除流动性、设置钱包地址等关键操作。该合约基于 Solidity 0.8 编写，继承自 OpenZeppelin 的 `ERC20` 和 `Ownable`，并集成了 Uniswap V2 路由器与交易对功能，具备自动税收分配与流动性添加机制。

## 二、合约核心特性概览
+ **代币总量**：1000 LMEME（18 位小数）
+ **买入税**：5%（用于营销 + 流动性）
+ **卖出税**：8%（用于营销 + 流动性）
+ **税费分配**：50% 营销钱包，50% 自动添加流动性
+ **交易冷却**：卖出后需等待 5 秒才能再次卖出
+ **最大单笔交易**：10% 总供应量（即 100 LMEME）
+ **自动做市**：卖出时自动将部分代币兑换为 ETH 并添加流动性
+ **免税地址**：合约自身、部署者、零地址、交易对、路由器
+ **测试**：支持在 Hardhat 环境中本地网络测试（含 console.log 调试）

## 三、部署前准备
**确保已安装以下工具**：

+ Node.js（v16+）
+ Hardhat
+ `@openzeppelin/contracts **v4版本**`
+ `@uniswap/v2-core`, `@uniswap/v2-periphery`
+ `hardhat/console.sol`

```markdown
npm install @openzeppelin/contracts@4.9.6 @uniswap/v2-core @uniswap/v2-periphery 
```

**获取测试网/主网参数**：

| **网络** | **Router 地址** | **WETH 地址** | **Factory 地址** |
| --- | --- | --- | --- |
| Ethereum (主网) | `0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D` | `0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2` | `0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f` |
| Goerli (测试网) | `0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D` | `0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6` | `0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f` |
| Sepolia (测试网) | `0xC532a74256D3Db42D0Bf7a0400fEFDbad72b3b9f` | `0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14` | `0x9B494d924787519a48333b993a56362E422098dD` |


## **四、部署合约**
合约部署脚本：scripts/deploy.js

部署命令：

```markdown
npx hardhat run scripts/deploy.js --network goerli
```

部署成功后会输出：代币合约地址，交易对地址（Pair）等信息

## **五、部署合约添加初始流动性**
**注意**：必须先向合约地址转入 ETH，再调用 addInitialLiquidity。

+ **发送 ETH 到合约地址**：使用钱包（如 MetaMask）向部署后的合约地址发送至少 0.1 ETH
+ **调用 **`**addInitialLiquidity**`** 函数**：

```markdown
// 示例调用（使用 ethers.js）
await token.addInitialLiquidity(
  ethers.utils.parseEther("500"), // 500 LMEME
  ethers.utils.parseEther("1")  // 1 ETH
);
```

## 六、代币交易规则
### **1. 买入代币（Buy）**
用户向交易对（Pair）发送 ETH，换取 LMEME。

+ **税率**：`5%`（从买入金额中扣除）
+ **去向**：
    - 2.5% → 营销钱包
    - 2.5% → 合约（后续用于流动性）
+ 无冷却时间限制。

### **2. 卖出代币（Sell）**
用户向交易对发送 LMEME，换取 ETH。

+ **税率**：`8%`
+ **去向**：
    - 4% → 营销钱包
    - 4% → 合约 → 自动兑换为 ETH 并添加流动性
+ **冷却时间**：卖出后需等待 `5 秒` 才能再次卖出
+ 若 swapAndLiquifyEnabled == false，则 4% 的流动性部分仍进入合约，但不执行添加。

## 七、关键功能操作
### 1. **设置最大交易额度**
仅限 **Owner** 调用。

```markdown
await token.setMaxTxAmount(ethers.utils.parseEther("200")); // 设置为 200 LMEME
```

### 2. **切换自动流动性功能**
```markdown
await token.updateSwapAndLiquifyEnabled(true); // 开启
await token.updateSwapAndLiquifyEnabled(false); // 关闭
```

## **八、免税与权限说明**
| **地址** | **是否免税** | **说明** |
| --- | --- | --- |
| `owner()` | ✅ | 部署者 |
| `address(this)` | ✅ | 合约自身 |
| `uniswapV2Pair` | ✅ | 交易对 |
| `uniswapV2Router` | ✅ | 路由器 |
| `address(0)` | ✅ | 零地址 |
| 其他用户 | ❌ | 正常征税 |


免税地址之间转账不触发税收逻辑。

## 九、监控与调试
### 1. **事件日志**
合约触发以下事件：

| **事件** | **说明** |
| --- | --- |
| `SwapAndLiquify` | 成功添加流动性 |
| `SwapAndLiquifyFailure` | 兑换失败 |
| `AddLiquidityFailure` | 添加流动性失败 |
| `MarketingWalletChanged` | 营销钱包变更（当前仅部署时触发） |


### 2. **Hardhat 调试**
使用 `console.log` 查看交易细节（仅限本地测试）：

```markdown
console.log("user sell");
console.log("tax:", tax);
```



## 部署及用户使用流程
### 阶段1：项目方部署和初始化
**1、部署**`**memetoken**`** 合约**

+ 项目方(Owner)通过部署脚本，将合约部署到链上，得到LMEME** 代币合约地址**（例如：`0x123...abc`）。
+ 构造函数执行：
    - 铸造 `1000LMEME` 到 `owner()` 钱包。
    - 设置 `marketingWallet = owner()`，`liquidityWallet = owner()`。
    - 设置免税地址：`owner()`、`this`、`uniswapV2Router`（占位，还未设置）。

**2、调用 **`**createPair(routerAddress)**`

+ 项目方调用：`memetoken.createPair(0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D);`
+ 合约内部：
    - 使用 Uniswap V2 Factory 创建 `LMEME/WETH` 交易对（Pair）。
    - 设置 `uniswapV2Pair = 0xdef...xyz`（新创建的 Pair 地址）。
    - 设置 `uniswapV2Router` 实例。
+ 注意：此时，交易对已存在，但**流动性为0**。

### 阶段2：添加初始流动性（市场启动）
**3、项目方添加流动性**

+ 项目方打开 Uniswap App，连接钱包。
+ 进入 “Pool” → “Add Liquidity”。
+ 选择 `LMEME`（输入合约地址）和 `WETH`。
+ 输入金额（如：500LMEME + 1 ETH）。
+ 点击 “ApproveLMEME” → 签名授权 Router 使用代币。
+ 点击 “Add Liquidity” → 签名第二笔交易。

底层发生了什么？

+ Uniswap Router 调用：

```plain
IUniswapV2Router02(router).addLiquidityETH{value: 1 ether}(
  address(memetoken), //LMEME 代币地址
  500e18,             //LMEME 数量：500，18位小数
  0, 0,               // slippage
  owner(),            // LP Token 接收者
  block.timestamp
);
```

+ Router 从项目方钱包扣除 500LMEME 和 1 ETH。
+ 注入 `LMEME/WETH` 交易对（Pair 合约）。
+ Pair 合约铸造 LP Token 并返回给项目方。
+ 此时池中有资金：1 ETH + 500LMEME ，流动性池建立，市场价格形成（如 1LMEME = 0.002 ETH，1 ETH = 500LMEME）。

### 阶段3：用户开始交易LMEME/WETH
现在，任何用户都可以通过 Uniswap 买卖LMEME。

#### 4、用户买入LMEME（用 ETH 换LMEME）
**（1）用户操作（前端）**

+ 用户连接钱包到 Uniswap。
+ 输入：`1 ETH` → 输出：`500LMEME`（含滑点）。
+ 点击 “Swap”。

**（2）底层调用（Router 执行）**

+ Uniswap Router 执行：

```plain
router.swapExactETHForTokens{value: 1 ether}(
    450e18,                     // 最少接收 450LMEME（滑点 10%）
    [WETH,LMEME],               // 路径：ETH → WETH →LMEME
    userAddress,                // 接收LMEME 的地址
    block.timestamp + 15
);
```

**（3）资金流与合约交互**

+ 用户的 1 ETH 被封装为 1 WETH，并**存入池中**。
+ Router 调用 `LMEME/WETH Pair` 合约进行 swap。
+ Pair 合约需要**从池中**划转LMEME 给用户。
+ 因此，Pair 合约调用：`memetoken.transfer(userAddress, amountOut);``**transfer**`** 是 ERC20 标准函数，它内部会调用 **`**_update(from, to, value)**`
+ **触发 memetoken**** 合约的 **`**_update()**`** 函数**：
    - `from = uniswapV2Pair`（因为是 Pair 发出LMEME）
    - `to = userAddress`
    - `value = amountOut`
+ 根据 memetoken 合约的代币税逻辑，此次操作为买入，税率为5%
+ `**_update()**`** 内部执行**：
    - 计算税费：`fee = value * 5 / 100`
    - 实际到账：`amountToTransfer = value - fee`
    - 调用 `super._update(from, to, amountToTransfer)` → 用户收到净额LMEME
    - 调用 `super._update(from, address(this), fee)` → 5% 税费转入 `memetoken ` 合约
+ 最终，用户成功买入LMEME， `memetoken ` 合约收到税费。

#### 5、用户卖出LMEME（用LMEME 换 ETH）
**（1）用户操作（前端）**

+ 用户输入：`100LMEME` → 输出：`~0.19 ETH`；点击 “Swap”。
+ 前提：用户已授权LMEME 合约：`memetoken.approve(routerAddress, 100e18);` 

**（2）底层调用（Router 执行）**

+ Uniswap Router 执行：

```plain
router.swapExactTokensForETH(
    100e18,
    0.18e18,              // 最少接收 ETH
    [LMEME, WETH],
    userAddress,
    block.timestamp + 15
);
```

**（3）资金流与合约交互**

+ Router 从用户钱包扣除 100LMEME。
+ 调用 `memetoken.transfer(uniswapV2Pair, 100e18)`
+ **触发 memetoken**** 合约的 **`**_update()**`** 函数**：
    - `from = userAddress`
    - `to = uniswapV2Pair`
    - `value = 100`
+ 根据 memetoken 合约的代币税逻辑，此次操作为卖出，税率为8%
+ `**update()**`** 内部执行**：
    - 检查冷却时间（如果是普通用户，且非免税）
    - 计算税费：`8% of 100LMEME = 8LMEME`
    - 用户实际转出 100LMEME：
        * `92LMEME` → 进入交易对（用于兑换 ETH）
        * `8LMEME` → 转入`memetoken`合约地址（税费）
+ 用户最终收到 ETH。

### 流程图
```plain
用户在 Uniswap 上点击 "Swap"
        ↓
Uniswap Router 调用 swap 函数
        ↓
Router 调用 memetoken.transfer(...) 转移LMEME
        ↓
ERC20.transfer() 内部调用 _update(from, to, value)
        ↓
你的 _update() 函数执行：
   - 判断买卖方向
   - 收取税费
   - 检查冷却
   - 触发 swapAndLiquify（当条件满足）
        ↓
交易完成，税费积累在合约中
        ↓
后续自动转化为流动性
```

## 说明
+ `memetoken` 合约只是一个 **ERC-20 代币合约**，它只管理LMEME 代币的发行、转账、税收等逻辑。
+ 想进行 LMEME/WETH 交易，就需要通过 `memetoken` 合约的 `createPair()` 函数初始化uniswap路由、创建LMEME/WETH 交易对，用户对于 LMEME/WETH 的买卖交易都是通过 LMEME/WETH 交易对进行的。
+ 交易对在进行交易的时候，会调用`memetoken` 合约的 `transfer()` 函数，从而触发 `_update()` 函数，实现`memetoken` 合约中的交易限制、代币税、流动性添加等业务。

## 关于 OpenZeppelin ^5.0.0中交易函数的覆写
在 v5 版本中：

+ `_transfer` 不再是 `virtual`，只是一个简单的检查函数，实际的转账逻辑委托给 `**_update**`
+ 代币合约中应该复写的函数应该是：

`function _update(address from, address to, uint256 value) internal virtual`

## 用户交易代币的流程图
**（1）基本设定**

+ **代币名称**：LMEME（MemeToken）
+ **税收机制**：
    - 买入税（Buy Tax）：5%（扣 LMEME）
    - 卖出税（Sell Tax）：8%（扣 LMEME）
+ **关键合约**：
    - `Router`：`IUniswapV2Router02`
    - `Pair`：`LMEME/WETH` 交易对合约
    - `LMEME`：代币合约（继承 `ERC20`，含税收逻辑）
    - `WETH`：Wrapped ETH 合约
+ **流动性池**：就是交易对合约`Pair`，流动性池的资金都储存在交易对合约中
+ **uniswap 路由**：就是`Router`，用户所有的交易都是通过`Router`进行的

**（2）流程图**

+ 用户买入 LMEME代币，用户用 ETH 换取 LMEME，触发 **5%的买入税**

![](https://cdn.nlark.com/yuque/0/2025/png/40797156/1761478010908-1d57babb-5455-4491-a9b5-bbcf794b5e4d.png)

+ **资产流动总结（买入，用户支付了 ETH，但税收是通过减少其获得的 LMEME 数量实现的）**

| **资产** | **用户** | **Pair（流动性池）** | **LMEME 合约（税收）** |
| --- | --- | --- | --- |
| ETH | - ETH | + ETH | 0 |
| LMEME | + 95% | - 100% | + 5% |




+ 用户卖出 LMEME，用户用 LMEME 换取 ETH，触发 **8%卖出税**

![](https://cdn.nlark.com/yuque/0/2025/png/40797156/1761478484636-bd43a27b-1b87-4f40-8e08-96483f2f75de.png)

+ **资产流动总结（卖出，所有资产变化在一个交易内原子完成）**

| **资产** | **用户** | **Pair（流动性池）** | **LMEME 合约（税收）** |
| --- | --- | --- | --- |
| LMEME | - 100% | + 92% | + 8% |
| ETH | + ETH | - ETH | 0 |


  
 

