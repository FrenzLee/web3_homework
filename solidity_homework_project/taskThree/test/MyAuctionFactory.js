const { ethers } = require("hardhat");
const { expect } = require("chai");

describe("MyAuctionFactory Test", function () {
    let MyAuction;
    let myAuction, myAuctionFactory, mynft;
    let nftAddress, myAuctionImplAddress;
    let owner, seller;

    const TOKEN_ID = 0;
    const MYNFT_NAME = "MyNFT";
    const MYNFT_SYMBOL = "MT";
    const TOKEN_URI = "https://example.com/my-nft.json";
    const DURATION = BigInt(5 * 24 * 60 * 60); // 5天
    const START_PRICE_ETH = ethers.parseEther("1"); // 1 ETH 起拍

    beforeEach(async function () { 
        //获取账户地址
        [owner, seller] = await ethers.getSigners();

        //部署拍卖逻辑合约
        MyAuction = await ethers.getContractFactory("MyAuction");
        myAuction = await MyAuction.deploy();
        await myAuction.waitForDeployment();
        myAuctionImplAddress = await myAuction.getAddress();

        //部署工厂合约
        const MyAuctionFactory = await ethers.getContractFactory("MyAuctionFactory");
        myAuctionFactory = await MyAuctionFactory.deploy(owner.address, myAuctionImplAddress);
        await myAuctionFactory.waitForDeployment();

        //部署NFT合约
        const MyNFT = await ethers.getContractFactory("MyNFT");
        mynft = await MyNFT.deploy(MYNFT_NAME, MYNFT_SYMBOL);
        await mynft.waitForDeployment();
        nftAddress = await mynft.getAddress();

        //铸造1个NFT给seller
        await mynft.mint(seller.address, TOKEN_URI);
        expect(await mynft.ownerOf(TOKEN_ID)).to.equal(seller.address);
    });

    it("Should deploy factory and set initial implementation", async function () {
        expect(await myAuctionFactory.auctionImplementation()).to.equal(myAuctionImplAddress);
        expect(await myAuctionFactory.implementationVersion()).to.equal(1);
        expect(await myAuctionFactory.owner()).to.equal(owner.address);
    });


    it("Should create a new auction proxy successfully", async function () {
        //预计算拍卖合约地址
        const predictedAddr = await myAuctionFactory.predictAuctionAddress(
            seller.address, nftAddress, TOKEN_ID, START_PRICE_ETH, DURATION, ethers.ZeroAddress);
        expect(predictedAddr).to.not.equal(ethers.ZeroAddress);        

        //卖家授权NFT
        await mynft.connect(seller).approve(predictedAddr, TOKEN_ID);

        //创建拍卖
        await expect(myAuctionFactory.connect(seller).createAuction(
            nftAddress,
            TOKEN_ID,
            START_PRICE_ETH,
            DURATION,
            ethers.ZeroAddress
        )).to.emit(myAuctionFactory, "AuctionDeployed")
          .withArgs(seller.address, predictedAddr, nftAddress, TOKEN_ID);

        //验证状态
        const auctionAddr = await myAuctionFactory.auctionsMap(nftAddress, TOKEN_ID);
        expect(auctionAddr).to.equal(predictedAddr);

        const allAuctions = await myAuctionFactory.getAllAuctions();
        expect(allAuctions.length).to.equal(1);
        expect(allAuctions[0]).to.equal(predictedAddr);
    });


    it("Should fail to create duplicate auction for same NFT", async function () {
        //预计算拍卖合约地址
        const predictedAddr = await myAuctionFactory.predictAuctionAddress(
            seller.address, nftAddress, TOKEN_ID, START_PRICE_ETH, DURATION, ethers.ZeroAddress);

        //卖家授权NFT
        await mynft.connect(seller).approve(predictedAddr, TOKEN_ID);

        // 第一次成功
        await myAuctionFactory.connect(seller).createAuction(
            nftAddress,
            TOKEN_ID,
            START_PRICE_ETH,
            DURATION,
            ethers.ZeroAddress
        )

        // 第二次失败
        await expect(
            myAuctionFactory.connect(seller).createAuction(
            nftAddress,
            TOKEN_ID,
            START_PRICE_ETH,
            DURATION,
            ethers.ZeroAddress
            )
        ).to.be.revertedWith("Auction already exists for this NFT");
    });


    it("Should initialize auction contract correctly", async function () {
        //预计算拍卖合约地址
        const predictedAddr = await myAuctionFactory.predictAuctionAddress(
            seller.address, nftAddress, TOKEN_ID, START_PRICE_ETH, DURATION, ethers.ZeroAddress);

        //卖家授权NFT
        await mynft.connect(seller).approve(predictedAddr, TOKEN_ID);

        // 创建合约
        await myAuctionFactory.connect(seller).createAuction(
            nftAddress,
            TOKEN_ID,
            START_PRICE_ETH,
            DURATION,
            ethers.ZeroAddress
        );

        // 获取代理合约实例
        const auctionProxy = MyAuction.attach(predictedAddr)
        const auctionInfo = await auctionProxy.getAuctionInfo();
        
        // 验证初始化参数
        expect(auctionInfo.seller).to.equal(seller.address);
        expect(auctionInfo.nftAddress).to.equal(nftAddress);
        expect(auctionInfo.tokenId).to.equal(TOKEN_ID);
        expect(auctionInfo.startPrice).to.equal(START_PRICE_ETH);
        expect(auctionInfo.endTime - auctionInfo.startTime).to.equal(DURATION);
        expect(auctionInfo.startTokenAddress).to.equal(ethers.ZeroAddress); 
    });

    it("Should upgrade auction implementation", async function () {

        const oldImpl = await myAuctionFactory.auctionImplementation()

        // 部署新版本拍卖逻辑
        const NewAuction = await ethers.getContractFactory("MyAuctionV2");
        const newImpl = await NewAuction.deploy();
        await newImpl.waitForDeployment();

        await expect(myAuctionFactory.upgradeImplementation(newImpl.target))
                .to.emit(myAuctionFactory, "AuctionImplementationUpgraded")
                .withArgs(2, oldImpl, newImpl.target);
        
        expect(await myAuctionFactory.auctionImplementation()).to.equal(newImpl.target);
        expect(await myAuctionFactory.implementationVersion()).to.equal(2);

        
        //预计算拍卖合约地址
        const predictedAddr = await myAuctionFactory.predictAuctionAddress(
            seller.address, nftAddress, TOKEN_ID, START_PRICE_ETH, DURATION, ethers.ZeroAddress);

        //卖家授权NFT
        await mynft.connect(seller).approve(predictedAddr, TOKEN_ID);

        // 创建合约
        await myAuctionFactory.connect(seller).createAuction(
            nftAddress,
            TOKEN_ID,
            START_PRICE_ETH,
            DURATION,
            ethers.ZeroAddress
        );

        // 获取代理合约实例
        const auctionProxy = NewAuction.attach(predictedAddr)
        const version = await auctionProxy.version();
        expect(version).to.equal("v2.0.0");
    });


});