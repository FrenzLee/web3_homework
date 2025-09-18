# 去中心化 NFT 拍卖平台

基于以太坊的去中心化 NFT 拍卖平台，支持多种代币支付，集成 Chainlink 价格预言机，并采用可升级代理模式。

## 目录

- [项目概述](#项目概述)
- [核心功能](#核心功能)
- [技术架构](#技术架构)
- [合约说明](#合约说明)
- [环境要求](#环境要求)
- [安装和配置](#安装和配置)
- [部署步骤](#部署步骤)
- [功能说明](#功能说明)
- [测试说明](#测试说明)
- [合约地址](#合约地址)
- [费用结构](#费用结构)
- [安全考虑](#安全考虑)

## 项目概述

本项目是一个功能完整的去中心化 NFT 拍卖平台，支持用户：
- 创建 NFT 拍卖
- 使用 ETH 或 ERC20 代币进行竞拍
- 基于实时汇率的美元计价系统
- 可升级的智能合约架构

## 核心功能

### 拍卖工厂 (MyAuctionFactory)
- **创建拍卖**：通过工厂模式创建标准化的拍卖合约
- **代理模式部署**：使用 UUPS 代理实现可升级性
- **实现合约管理**：支持升级拍卖合约实现
- **拍卖记录**：维护所有已创建拍卖的索引

### 拍卖合约 (MyAuction)
- **币种支付**：支持 ETH 和 USDC 等 ERC20 代币
- **实时价格转换**：集成 Chainlink 价格预言机
- **美元计价**：所有出价以美元为基准进行比较
- **安全机制**：防重入攻击、权限控制等

### MyNFT 合约
- **标准 ERC721**：完全兼容的 NFT 实现
- **铸造**：支持铸造 NFT
- **元数据存储**：支持自定义 URI 和元数据

## 合约说明

### MyAuctionFactory.sol
拍卖工厂合约，负责创建和管理拍卖实例。

**主要功能：**
- `createAuction()` - 创建新的拍卖
- `getAllAuctions()` - 获取所有拍卖列表
- `upgradeImplementation()` - 升级拍卖合约实现
- `predictAuctionAddress()` - 预计算代理合约地址

### MyAuction.sol / MyAuctionV2.sol
拍卖逻辑合约，采用可升级代理模式。

**主要功能：**
- `bidByETH()` - 使用 ETH 竞拍
- `bidByToken()` - 使用 ERC20 代币竞拍
- `endAuction()` - 结束拍卖
- `claimNFT()` - 领取 NFT
- `claimPayment()` - 卖家领取资金
- `claimNFTtoSeller()` - 无人竞拍卖家取回 NFT
- `withdrawStuckETH()` - 管理员提取意外收到的 ETH

### MyNFT.sol
标准的 ERC721 NFT 合约。

**主要功能：**
- `mint()` - 铸造单个 NFT

## 🔧 环境要求

- Node.js >= 16.0.0
- npm 或 yarn
- Git

## 安装和配置

### 1. 克隆项目
```bash
git clone https://github.com/FrenzLee/web3_homework.git
cd solidity_homework_project/taskThree
```

### 2. 安装依赖
```bash
npm install
```

## 部署步骤

### 1. 编译合约
```bash
npx hardhat compile
```

### 2. 运行测试
```bash
# 运行所有测试
npx hardhat test

# 运行特定测试
npx hardhat test test/MyAuction.js --network localhost
npx hardhat test test/MyAuctionFactory.js --network localhost
npx hardhat test test/MyNFT.js --network localhost

# 生成测试覆盖率报告
npx hardhat coverage
```

### 3. 部署到本地测试网
```bash
npx hardhat run scripts/deploy.js --network localhost
```

部署完成后，会在项目根目录生成 `deployment.json` 文件，包含所有合约地址。

### 4. 升级拍卖合约（可选）
如果需要升级拍卖合约实现：
```bash
npx hardhat run scripts/upgrade_auction.js --network localhost
```
部署完成后，会在项目根目录生成 `deployment-upgrade.json` 文件，包含所有合约地址。

## 功能说明

### 创建 NFT 拍卖

调用 `createAuction()` 函数创建新的拍卖。注意：卖家创建拍卖前，需要确保 MyAuctionFactory 合约已获得转移该 NFT 的授权。

```solidity
// 创建拍卖
myAuctionFactory.createAuction(
    seller,           // 卖家地址
    nftAddress,       // NFT 合约地址
    tokenId,          // NFT Token ID
    startingPrice,    // 起拍价格（ETH，18位小数）
    duration,         // 拍卖时长（秒）
    ethPriceFeed      // ETH 价格预言机地址
);
```

**工作流程：**
1. 创建合约工厂，并把 NFT 授权给 MyAuctionFactory
2. 创建拍卖代理合约并初始化
3. 竞拍者进行竞拍，并使用 ETH 或 ERC20 代币支付
4. 拍卖结束后，买家领取 NFT，卖家领取资金

### 参与竞拍

#### 使用 ETH 竞拍
```solidity
// 直接发送 ETH
myAuction.bidByETH({value: bidAmount});
```

#### 使用 ERC20 代币竞拍
```solidity
// 1. 授权代币转移
myerc20.approve(auctionAddress, bidAmount);

// 2. 进行竞拍
myAuction.bidByToken(tokenAddress, bidAmount);
```

### 结束拍卖和领取

```solidity
// 1. 结束拍卖
myAuction.endAuction();

// 2. 获胜者领取 NFT
myAuction.claimNFT();

// 3. 卖家领取付款
myAuction.claimPayment();
```

### Sepolia 测试网预言机地址
- **ETH/USD**: `0x694AA1769357215DE4FAC081bf1f309aDC325306`
- **USDC/USD**: `0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E`


## 安全考虑

### 已实现的安全措施
- **重入攻击防护**：使用 ReentrancyGuard
- **权限控制**：基于 OpenZeppelin Ownable
- **安全的代币转移**：使用 SafeERC20
- **输入验证**：全面的参数检查
- **时间锁定**：拍卖时间控制
- **价格验证**：预言机数据验证

## 开发工具

- **Hardhat**：开发环境和测试框架
- **OpenZeppelin**：安全的智能合约库
- **Chainlink**：去中心化预言机网络
- **Ethers.js**：以太坊交互库

**⚠️ 免责声明：本项目仅用于学习和研究目的。在生产环境中使用前，请进行充分的测试和安全审计。**
