const hre = require("hardhat");
const { ethers } = hre;
const fs = require('fs');
const path = require('path');

async function main() {
    const [deployer] = await ethers.getSigners();
    const network = await ethers.provider.getNetwork();
    const networkName = network.name === "unknown" ? "localhost" : network.name;

    console.log(`Upgrading MyAuction contract on network: ${networkName}`);
    console.log("Upgrading with account:", deployer.address);

    //读取之前的部署信息
    const deploymentPath = path.join(__dirname, "..", "deployments", networkName, "deployment.json")
    if (!fs.existsSync(deploymentPath)) {
        throw new Error(`Deployment file not found: ${deploymentPath}. Please deploy first.`);
    }

    const deploymentData = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
    const factoryAddress = deploymentData.contracts.MyAuctionFactory;
    console.log("Using MyAuctionFactory at:", factoryAddress);

    //部署 MyAuctionV2 实现合约
    const MyAuctionV2 = await ethers.getContractFactory("MyAuctionV2");
    const myAuctionV2Impl = await MyAuctionV2.deploy();
    await myAuctionV2Impl.waitForDeployment();
    const newImplAddress = await myAuctionV2Impl.getAddress();
    console.log("MyAuctionV2 Implementation deployed to:", newImplAddress);

    //获取 MyAuctionFactory 实例
    const MyAuctionFactory = await ethers.getContractFactory("MyAuctionFactory");
    const myAuctionFactory = MyAuctionFactory.attach(factoryAddress);

    //升级
    const tx = await myAuctionFactory.upgradeImplementation(newImplAddress);
    await tx.wait();
    console.log("Upgrade transaction confirmed:", tx.hash);

    //验证升级结果
    const currentImpl = await myAuctionFactory.auctionImplementation();
    const version = await myAuctionFactory.implementationVersion();
    console.log("Upgrade successful!");
    console.log("New Auction Implementation:", currentImpl);
    console.log("New Implementation Version:", version.toString());

    //构建部署信息
    deploymentData.contracts.MyAuctionV2Impl = newImplAddress;
    deploymentData.contracts.auctionImplInFactory = currentImpl;
    deploymentData.upgrade = {
        timestamp: new Date().toISOString(),
        from: deploymentData.contracts.MyAuctionImpl,
        to: newImplAddress,
        version: version.toString(),
        txHash: tx.hash
    };

    const outputPath = path.join(__dirname, "..", "deployments", networkName, "deployment-upgrade.json");
    fs.writeFileSync(outputPath, JSON.stringify(deploymentData, null, 2));
    console.log(`Deployment info updated at: ${outputPath}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("Upgrade failed:", error);
    process.exit(1);
  });