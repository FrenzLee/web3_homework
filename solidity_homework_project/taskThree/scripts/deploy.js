const hre = require("hardhat");
const { ethers } = hre;
const fs = require('fs');
const path = require('path');

async function main() { 
    const [deployer] = await ethers.getSigners();
    const network = await ethers.provider.getNetwork();
    const networkName = network.name === "unknown" ? "localhost" : network.name;

    console.log(`Deploying to network: ${networkName}`);
    console.log("Deploying contracts with account:", deployer.address);
    //console.log("Account balance:", (await deployer.getBalance()).toString());

    // 部署 MyAuction 逻辑合约
    const MyAuction = await ethers.getContractFactory("MyAuction");
    const myAuctionImpl = await MyAuction.deploy();
    await myAuctionImpl.waitForDeployment();
    const implAddr = await myAuctionImpl.getAddress();
    console.log("MyAuction Implementation deployed to:", implAddr);

    // 部署 MyAuctionFactory 的实例
    const MyAuctionFactory = await ethers.getContractFactory("MyAuctionFactory");
    const myAuctionFactoryImpl = await MyAuctionFactory.deploy(deployer.address, implAddr);
    await myAuctionFactoryImpl.waitForDeployment();
    const factoryImplAddr = await myAuctionFactoryImpl.getAddress();
    console.log("MyAuctionFactory deployed to:", factoryImplAddr);

    // 读取工厂中存储的实现地址
    const auctionImplAddr = await myAuctionFactoryImpl.auctionImplementation();

    // 构建部署信息
    const deploymentData = {
        network: networkName,
        deployer: deployer.address,
        timestamp: new Date().toISOString(),
        contracts: {
            MyAuctionImpl: implAddr,
            MyAuctionFactory: factoryImplAddr,
            auctionImplInFactory: auctionImplAddr,
        },
        rpcUrl: hre.network.config.url,
    };
    
    // 写入 JSON 文件
    const outputDir = path.join(__dirname, "..", "deployments", networkName);
    if (!fs.existsSync(outputDir)) {
        fs.mkdirSync(outputDir, { recursive: true });
    }
    const outputPath = path.join(outputDir, "deployment.json");
    fs.writeFileSync(outputPath, JSON.stringify(deploymentData, null, 2));

    console.log(`Deployment info saved to ${outputPath}`);

    console.log("All done!");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

