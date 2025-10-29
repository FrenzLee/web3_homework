const { ethers, upgrades } = require('hardhat');
const { expect } = require('chai');

describe("MetaNode Test", function () {
    let metaNodeStake, metaNodeToken, boyToken;
    let metaNodeStakeAddr, metaNodeTokenAddr, boyTokenAddr;
    let owner, user1, user2;
    const ETH_PID = 0;

    this.beforeEach(async function() {
        [owner, user1, user2] = await ethers.getSigners();

        //部署MetaNodeToken，奖励代币
        const MetaNodeToken = await ethers.getContractFactory("MetaNodeToken", owner);
        metaNodeToken = await MetaNodeToken.deploy();
        await metaNodeToken.waitForDeployment();
        metaNodeTokenAddr = await metaNodeToken.getAddress();

        //部署BoyToken，质押代币
        const BoyToken = await ethers.getContractFactory("BoyToken", owner);
        boyToken = await BoyToken.deploy();
        await boyToken.waitForDeployment();
        boyTokenAddr = await boyToken.getAddress();

        //部署MetaNodeStake，质押合约
        const MetaNodeStake = await ethers.getContractFactory("MetaNodeStake", owner);
        //使用 upgrades 部署 UUPS 代理合约
        const startBlock = 0; //模拟从 0 开始
        const endBlock = 1000;
        const metaNodePerBlock = ethers.parseEther("10"); //每个区块奖励10个MetaNodeToken代币
        metaNodeStake = await upgrades.deployProxy(
            MetaNodeStake, 
            [metaNodeTokenAddr, startBlock, endBlock, metaNodePerBlock],
            { initializer: "initialize" }
        );
        await metaNodeStake.waitForDeployment();
        metaNodeStakeAddr = await metaNodeStake.getAddress();

        console.log("beforeEach is OK ! ");
        console.log("MetaNodeStake deployed to:", metaNodeStakeAddr);
        console.log("MetaNodeToken deployed to:", metaNodeTokenAddr);
        console.log("BoyToken deployed to:", boyTokenAddr);
    });

    it("Should deploy and set correct parameters", async function () {
        expect(await metaNodeStake.MetaNode()).to.equal(await metaNodeToken.getAddress());
        expect(await metaNodeStake.startBlock()).to.equal(0);
        expect(await metaNodeStake.endBlock()).to.equal(1000);
        expect(await metaNodeStake.MetaNodePerBlock()).to.equal(ethers.parseEther("10"));
    });
	
	
	it("Should add pool correctly", async function () {
        await metaNodeStake.addPool(ethers.ZeroAddress, 100, ethers.parseEther("0.01"), 10, false);

        expect(await metaNodeStake.poolLength()).to.equal(1);

        const pool = await metaNodeStake.pool(ETH_PID);
        expect(pool.stTokenAddress).to.equal(ethers.ZeroAddress);
        expect(pool.poolWeight).to.equal(100);
        expect(pool.minDepositAmount).to.equal(ethers.parseEther("0.01"));
        expect(pool.unstakeLockedBlocks).to.equal(10);
    });
	
	
	it("Should allow user to deposit ETH and earn rewards", async function () {
        await metaNodeStake.addPool(ethers.ZeroAddress, 100, ethers.parseEther("0.01"), 10, false);

        const depositAmount = ethers.parseEther("1");
        await metaNodeStake.connect(user1).depositETH({ value: depositAmount });

        const user = await metaNodeStake.user(ETH_PID, user1.address);
        expect(user.stAmount).to.equal(depositAmount);

        // 快进区块
        await network.provider.send("evm_mine");

        const pending = await metaNodeStake.pendingMetaNode(ETH_PID, user1.address);
        console.log("pendingMetaNode:", pending);
        expect(pending).to.be.greaterThan(0);

    });
	
    it("Should allow user to claim rewards", async function () {
        // 将 owner 的 MetaNodeToken 转给 metaNodeStake 合约
        const transferAmount = ethers.parseEther("10000000"); // 假设转 10000000 个代币
        await metaNodeToken.connect(owner).transfer(metaNodeStakeAddr, transferAmount);

        await metaNodeStake.addPool(ethers.ZeroAddress, 100, ethers.parseEther("0.01"), 10, false);
        await metaNodeStake.connect(user1).depositETH({ value: ethers.parseEther("1") });

        const pendingBefore = await metaNodeStake.pendingMetaNode(ETH_PID, user1.address);
        console.log("pendingBefore:", pendingBefore);

        // 快进区块
        for (let i = 0; i < 5; i++) await network.provider.send("evm_mine");
        // 领取奖励
        await metaNodeStake.connect(user1).claim(ETH_PID);
        
        const balanceAfterClaim = await metaNodeToken.balanceOf(user1.address);
        console.log("balanceAfterClaim:", balanceAfterClaim);
        expect(balanceAfterClaim).to.greaterThan(pendingBefore);
    });
	
    it("Should allow ERC20 token deposit and claim", async function () {
        // 将 owner 的 MetaNodeToken 转给 metaNodeStake 合约
        const transferAmount = ethers.parseEther("10000000"); // 假设转 10000000 个代币
        await metaNodeToken.connect(owner).transfer(metaNodeStakeAddr, transferAmount);

        // 先添加ETH质押池
        await metaNodeStake.addPool(ethers.ZeroAddress, 100, ethers.parseEther("0.01"), 10, false);

        // 授权
        await boyToken.connect(user1).approve(metaNodeStake, ethers.parseEther("100"));

        // 添加 BoyToken 质押池
        await metaNodeStake.addPool(await boyToken.getAddress(), 200, 1, 5, false);

        // 质押 BoyToken
        await boyToken.transfer(user1.address, ethers.parseEther("100"));
        const depositAmount = ethers.parseEther("50");
        await metaNodeStake.connect(user1).deposit(1, depositAmount);

        const user = await metaNodeStake.user(1, user1.address);
        console.log("user.stAmount:", user.stAmount);
        expect(user.stAmount).to.equal(depositAmount);

        const pendingBefore = await metaNodeStake.pendingMetaNode(ETH_PID, user1.address);
        console.log("pendingBefore:", pendingBefore);

        // 快进区块
        for (let i = 0; i < 5; i++) await network.provider.send("evm_mine");

        // 领取奖励
        await metaNodeStake.connect(user1).claim(1);
        const balanceAfterClaim = await metaNodeToken.balanceOf(user1.address);
        console.log("balanceAfterClaim:", balanceAfterClaim);
        expect(balanceAfterClaim).to.greaterThan(pendingBefore);

    });
	
    it("Should handle unstake and withdraw flow", async function () {
        // 将 owner 的 MetaNodeToken 转给 metaNodeStake 合约
        const transferAmount = ethers.parseEther("10000000"); // 假设转 10000000 个代币
        await metaNodeToken.connect(owner).transfer(metaNodeStakeAddr, transferAmount);

        // 先添加ETH质押池
        await metaNodeStake.addPool(ethers.ZeroAddress, 100, ethers.parseEther("0.01"), 3, false);

        const depositAmount = ethers.parseEther("2");
        await metaNodeStake.connect(user1).depositETH({ value: depositAmount });

        // 快进
        await network.provider.send("evm_mine");

        // 申请解质押 1 ETH
        await metaNodeStake.connect(user1).unstake(ETH_PID, ethers.parseEther("1"));

        // 检查请求
        const [requestAmount, pendingWithdraw] = await metaNodeStake.withdrawAmount(ETH_PID, user1.address);
        expect(requestAmount).to.equal(ethers.parseEther("1"));
        expect(pendingWithdraw).to.equal(0); // 尚未解锁

        // 快进超过解锁区块
        for (let i = 0; i < 6; i++) await network.provider.send("evm_mine");

        const [, newPendingWithdraw] = await metaNodeStake.withdrawAmount(ETH_PID, user1.address);
        console.log("newPendingWithdraw:", newPendingWithdraw);
        expect(newPendingWithdraw).to.equal(ethers.parseEther("1"));
        
        // 提现
        await expect(() => metaNodeStake.connect(user1).withdraw(ETH_PID))
            .to.changeEtherBalance(user1, ethers.parseEther("1"));
    });

    it("Should allow admin to pause claim and withdraw", async function () {
        await metaNodeStake.pauseClaim();
        expect(await metaNodeStake.claimPaused()).to.be.true;

        await metaNodeStake.pauseWithdraw();
        expect(await metaNodeStake.withdrawPaused()).to.be.true;

        await metaNodeStake.unpauseClaim();
        await metaNodeStake.unpauseWithdraw();

        expect(await metaNodeStake.claimPaused()).to.be.false;
        expect(await metaNodeStake.withdrawPaused()).to.be.false;
    });
	
    it("Should not allow deposit below min amount", async function () {
        await metaNodeStake.addPool(ethers.ZeroAddress, 100, ethers.parseEther("1"), 10, false);
        await expect(
            metaNodeStake.connect(user1).depositETH({ value: ethers.parseEther("0.5") })
        ).to.be.revertedWith("deposit amount is too small");
    });

});