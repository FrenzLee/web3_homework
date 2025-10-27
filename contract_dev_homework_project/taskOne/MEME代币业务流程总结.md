## 一、流程整理
### 阶段1：项目方部署和初始化
**1、部署**`memetoken`合约

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

## 二、用户交易代币的流程图
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


## 三、关于 OpenZeppelin ^5.0.0中交易函数的覆写
在 v5 版本中：

+ `_transfer` 不再是 `virtual`，只是一个简单的检查函数，实际的转账逻辑委托给 `**_update**`
+ 代币合约中应该复写的函数应该是：

`function _update(address from, address to, uint256 value) internal virtual`  
 

