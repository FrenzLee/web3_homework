const { ethers, upgrades } = require("hardhat");
const { expect } = require("chai");

describe("MyAuction Test", function () { 
    let auctionProxy;
    let mynft, myerc20;
    let ethToUsdPriceFeed, erc20ToUsdPriceFeed; 
    let owner, seller, bidder1, bidder2, bidder3;

    const MYNFT_NAME = "MyNFT";
    const MYNFT_SYMBOL = "MT";
    const TOKEN_URI = "https://example.com/my-nft.json";
    const MYERC20_NAME = "MyERC20";
    const MYERC20_SYMBOL = "ME";

    const DURATION = BigInt(5 * 24 * 60 * 60); // 5天
    const ETH_PRICE = ethers.parseUnits("2000", 8); // $2000 in 8 decimals
    const ERC20_PRICE = ethers.parseUnits("500", 8); // $500 in 8 decimals
    const START_PRICE_ETH = ethers.parseEther("1"); // 1 ETH 起拍


    beforeEach(async () => { 
        //获取账户地址
        [owner, seller, bidder1, bidder2, bidder3] = await ethers.getSigners();

        //owner部署基础合约
        const MYNFT = await ethers.getContractFactory("MyNFT");
        mynft = await MYNFT.connect(owner).deploy(MYNFT_NAME, MYNFT_SYMBOL);
        await mynft.waitForDeployment();

        const MYERC20 = await ethers.getContractFactory("MyERC20");
        myerc20 = await MYERC20.connect(owner).deploy(MYERC20_NAME, MYERC20_SYMBOL, 10000);
        await myerc20.waitForDeployment();

        const ETHToUsdPriceFeed = await ethers.getContractFactory("EthToUsdPriceFeed");
        ethToUsdPriceFeed = await ETHToUsdPriceFeed.connect(owner).deploy(8, ETH_PRICE);
        await ethToUsdPriceFeed.waitForDeployment();

        const ERC20ToUsdPriceFeed = await ethers.getContractFactory("Erc20ToUsdPriceFeed");
        erc20ToUsdPriceFeed = await ERC20ToUsdPriceFeed.connect(owner).deploy(8, ERC20_PRICE);
        await erc20ToUsdPriceFeed.waitForDeployment();

        const MyAuction = await ethers.getContractFactory("MyAuction");
        auctionProxy = await upgrades.deployProxy(MyAuction, [], {initializer: false,});
        await auctionProxy.waitForDeployment();


    });
    

    describe("deployment", function () {
        it("Should be deployed correctly", async function () { 
            expect(await mynft.name()).to.equal(MYNFT_NAME);
            expect(await mynft.symbol()).to.equal(MYNFT_SYMBOL);
            expect(await mynft.owner()).to.equal(owner.address);
            expect(await myerc20.name()).to.equal(MYERC20_NAME);
            expect(await myerc20.symbol()).to.equal(MYERC20_SYMBOL);
            expect(await myerc20.owner()).to.equal(owner.address);

            const [ ,ethToUsdPrice ] = await ethToUsdPriceFeed.latestRoundData();
            const [ ,erc20ToUsdPrice ] = await erc20ToUsdPriceFeed.latestRoundData();
            expect(await ethToUsdPriceFeed.owner()).to.equal(owner.address);
            expect(await ethToUsdPriceFeed.decimals()).to.equal(8);
            expect(await ethToUsdPriceFeed.description()).to.equal("EthToUsdPriceFeed");
            expect(await ethToUsdPriceFeed.version()).to.equal(1);
            expect(ethToUsdPrice).to.equal(ETH_PRICE);
            expect(await erc20ToUsdPriceFeed.owner()).to.equal(owner.address);
            expect(await erc20ToUsdPriceFeed.decimals()).to.equal(8);
            expect(await erc20ToUsdPriceFeed.description()).to.equal("Erc20ToUsdPriceFeed");
            expect(await erc20ToUsdPriceFeed.version()).to.equal(1);
            expect(erc20ToUsdPrice).to.equal(ERC20_PRICE);            
        });
    });

    
    describe("Initialization", function () { 
        const tokenId = 0;

        beforeEach(async () => {
            //卖家创建NFT，第一次tokenId=0
            await mynft.mint(seller.address, TOKEN_URI);
            //卖家授权
            await mynft.connect(seller).approve(auctionProxy, tokenId);
        });

        it("Should be initialized correctly", async function () { 
            //mynft.target和await mynft.getAddress()得到的都是mynft的合约地址
            //await mynft.ownerOf(0)得到的是tokenId=0的owner地址
            //创建合约
            await auctionProxy.initialize(
                seller.address,
                mynft.target,
                tokenId,
                START_PRICE_ETH,
                DURATION,
                ethers.ZeroAddress
            );

            //设置预言机
            await auctionProxy.setPriceFeeds(ethers.ZeroAddress, ethToUsdPriceFeed.target);
            await auctionProxy.setPriceFeeds(myerc20.target, erc20ToUsdPriceFeed.target);

            const auctionInfo = await auctionProxy.getAuctionInfo();
            
            //验证合约初始化信息
            expect(auctionInfo.seller).to.equal(seller.address);
            expect(auctionInfo.nftAddress).to.equal(mynft.target);
            expect(auctionInfo.tokenId).to.equal(tokenId);
            expect(auctionInfo.startPrice).to.equal(START_PRICE_ETH);
            expect(auctionInfo.highestPrice).to.equal(START_PRICE_ETH);
            expect(auctionInfo.highestBidder).to.equal(ethers.ZeroAddress);
            expect(auctionInfo.highestTokenAddress).to.equal(ethers.ZeroAddress);
            expect(auctionInfo.startTokenAddress).to.equal(ethers.ZeroAddress);
            expect(auctionInfo.startTime).to.be.gt(0);
            expect(auctionInfo.endTime).to.equal(BigInt(auctionInfo.startTime) + DURATION);
            expect(auctionInfo.isEnd).to.be.false;
            expect(auctionInfo.nftClaimed).to.be.false;
            expect(auctionInfo.paymentClaimed).to.be.false;

            //验证nft授权
            expect(await mynft.getApproved(tokenId)).to.equal(await auctionProxy.getAddress());
            //此时NFT归属seller，因为只是给拍卖合约授权了，没有转移NFT
            expect(await mynft.ownerOf(tokenId)).to.equal(seller.address);

            //校验价格预言机，格式化返回的价格
            const ethToUsd = await auctionProxy.getPriceUSD(ethers.ZeroAddress, 1);
            const erc20ToUsd = await auctionProxy.getPriceUSD(myerc20.target, 1);
            // 转为 8 位小数的整数表示
            const to8Decimals = (num) => num * BigInt(10 ** 8);

            expect(to8Decimals(ethToUsd)).to.equal(ETH_PRICE);
            expect(to8Decimals(erc20ToUsd)).to.equal(ERC20_PRICE);

        });

    }); 

    
    describe("Bidding with ETH", function () { 
        const tokenId = 0;

        beforeEach(async () => {
            //卖家创建NFT，第一次tokenId=0
            await mynft.mint(seller.address, TOKEN_URI);
            //卖家授权
            await mynft.connect(seller).approve(auctionProxy, tokenId);

            //创建合约
            await auctionProxy.initialize(
                seller.address,
                mynft.target,
                tokenId,
                START_PRICE_ETH,
                DURATION,
                ethers.ZeroAddress
            );

            //设置预言机
            await auctionProxy.setPriceFeeds(ethers.ZeroAddress, ethToUsdPriceFeed.target);
            await auctionProxy.setPriceFeeds(myerc20.target, erc20ToUsdPriceFeed.target);
        });

        
        it("Should accept highest bid in ETH", async function () { 
            const bidAmount = ethers.parseEther("1.5");

            //触发HighestBidIncreased事件校验
            expect(await auctionProxy.connect(bidder1).bidByETH({value: bidAmount}))
                .to.emit(auctionProxy, "HighestBidIncreased")
                .withArgs(bidder1.address, tokenId, ethers.ZeroAddress, bidAmount);
            
            const auctionInfo = await auctionProxy.getAuctionInfo();
            expect(auctionInfo.highestBidder).to.equal(bidder1.address);
            expect(auctionInfo.highestPrice).to.equal(bidAmount);
            expect(auctionInfo.highestTokenAddress).to.equal(ethers.ZeroAddress);
        });  

        
        it("Should refund previous bidder", async function () {
            const firstBid = ethers.parseEther("1.2");
            const secondBid = ethers.parseEther("1.5");

            //bidder1竞拍前的余额
            const bidder1BlanceBefore = await ethers.provider.getBalance(bidder1.address);
            const bidder2BlanceBefore = await ethers.provider.getBalance(bidder2.address);
            const auctionBlanceBefore = await ethers.provider.getBalance(auctionProxy.target);

            //bidder1竞拍
            await auctionProxy.connect(bidder1).bidByETH({value: firstBid});
            //bidder2竞拍，价格更高，bidder1的钱应该退回
            await auctionProxy.connect(bidder2).bidByETH({value: secondBid});

            //bidder1被退回后的余额
            const bidder1BlanceFinal = await ethers.provider.getBalance(bidder1.address);
            const bidder2BlanceFinal = await ethers.provider.getBalance(bidder2.address);
            const auctionBlanceFinal = await ethers.provider.getBalance(auctionProxy.target);

            //验证bidder2是最高出价者
            const auctionInfo = await auctionProxy.getAuctionInfo();
            expect(auctionInfo.highestBidder).to.equal(bidder2.address);
            expect(auctionInfo.highestPrice).to.equal(secondBid);
            expect(auctionInfo.highestTokenAddress).to.equal(ethers.ZeroAddress);

            //验证bidder1退回后的余额，因为会扣除gas费，所以前后的余额应该有些差距
            expect(bidder1BlanceBefore).to.be.closeTo(bidder1BlanceFinal, ethers.parseEther("0.01"));
            // console.log("bidder1BlanceBefore:", bidder1BlanceBefore);
            // console.log("bidder1BlanceFinal:", bidder1BlanceFinal);
            // console.log("bidder2BlanceBefore:", bidder2BlanceBefore);
            // console.log("bidder2BlanceFinal:", bidder2BlanceFinal);
            // console.log("auctionBlanceBefore:", auctionBlanceBefore);
            // console.log("auctionBlanceFinal:", auctionBlanceFinal);

        }); 

        
        it("Should reject bid lower than start price", async function () { 
            const lowBid = ethers.parseEther("0.5");
            
            await expect(auctionProxy.connect(bidder1).bidByETH({value: lowBid}))
                .to.be.revertedWith("Bid amount must be higher than the start price");
        });

        it("Should reject bid not higher than current highest", async function () { 
            const bid1 = ethers.parseEther("2");
            const bid2 = ethers.parseEther("1.8");

            await auctionProxy.connect(bidder1).bidByETH({value: bid1});
            
            await expect(auctionProxy.connect(bidder2).bidByETH({value: bid2}))
                .to.be.revertedWith("Bid amount must be higher than the current highest bid");
        });  

    }); 

    
    describe("Bidding with ERC20", function () { 
        const tokenId = 0;

        beforeEach(async () => {
            //卖家创建NFT，第一次tokenId=0
            await mynft.mint(seller.address, TOKEN_URI);
            //卖家授权
            await mynft.connect(seller).approve(auctionProxy, tokenId);

            //创建合约
            await auctionProxy.initialize(
                seller.address,
                mynft.target,
                tokenId,
                START_PRICE_ETH,
                DURATION,
                ethers.ZeroAddress
            );

            //设置预言机
            await auctionProxy.setPriceFeeds(ethers.ZeroAddress, ethToUsdPriceFeed.target);
            await auctionProxy.setPriceFeeds(myerc20.target, erc20ToUsdPriceFeed.target);

            //给竞拍者发放代币并授权
            //如果不写ethers.parseUnits("100", 18)，直接写100，代表100 * 10**-18 wei
            await myerc20.transfer(bidder1.address, ethers.parseUnits("100", 18));
            await myerc20.transfer(bidder2.address, ethers.parseUnits("100", 18));
            await myerc20.connect(bidder1).approve(auctionProxy, ethers.parseUnits("100", 18));
            await myerc20.connect(bidder2).approve(auctionProxy, ethers.parseUnits("100", 18));
        });

        
        it("Should accept highest bid in ERC20", async function () { 
            
            //查看代币供应数量
            //const totalSupply = await myerc20.totalSupply();
            //console.log("Total Supply (raw):", totalSupply.toString());
            // 如果 decimals=18，应该是 10000.0
            //console.log("Total Supply (formatted):", ethers.formatUnits(totalSupply, 18));
            

            const lowBidAmount = ethers.parseUnits("2", 18);
            const middleBidAmount = ethers.parseUnits("4", 18);
            const highBidAmount = ethers.parseUnits("10", 18);
            
            //const aaa = await auctionProxy.getPriceUSD(ethers.ZeroAddress, START_PRICE_ETH);
            //const bbb = await auctionProxy.getPriceUSD(myerc20.target, highBidAmount);
            //console.log("aaa:", aaa);
            //console.log("bbb:", bbb);
            //console.log("1 ETH 价值 USD:", ethers.formatUnits(aaa, 18)); // "2000.0"
            //console.log("10 MyERC20 价值 USD:", ethers.formatUnits(bbb, 18)); // "5000.0"        

            //bidder1出价低于起拍价
            await expect(auctionProxy.connect(bidder1).bidByToken(myerc20.target, lowBidAmount))
                .to.be.revertedWith("Bid amount must be higher than the start price");
            
            
            //bidder2出价最高
            await auctionProxy.connect(bidder2).bidByToken(myerc20.target, highBidAmount);
            const auctionInfo = await auctionProxy.getAuctionInfo();
            expect(auctionInfo.highestBidder).to.equal(bidder2.address);
            expect(auctionInfo.highestPrice).to.equal(highBidAmount);
            expect(auctionInfo.highestTokenAddress).to.equal(myerc20.target);
            

            //bidder1出价等于起拍价，低于bidder2出价
            await expect(auctionProxy.connect(bidder1).bidByToken(myerc20.target, middleBidAmount))
                .to.be.revertedWith("Bid amount must be higher than the current highest bid");

        });  

        it("Should refund previous bidder", async function () {
            const firstBid = ethers.parseUnits("10", 18);
            const secondBid = ethers.parseUnits("15", 18);

            //bidder1竞拍前的余额
            const bidder1BlanceBefore = await myerc20.balanceOf(bidder1.address);
            const bidder2BlanceBefore = await myerc20.balanceOf(bidder2.address);
            const auctionBlanceBefore = await myerc20.balanceOf(auctionProxy.target);

            //bidder1竞拍
            await auctionProxy.connect(bidder1).bidByToken(myerc20.target, firstBid);
            const auctionInfoBid1 = await auctionProxy.getAuctionInfo();
            expect(auctionInfoBid1.highestPrice).to.equal(firstBid);

            //bidder2竞拍，价格更高，bidder1的钱应该退回
            await auctionProxy.connect(bidder2).bidByToken(myerc20.target, secondBid);

            //bidder1被退回后的余额
            const bidder1BlanceFinal = await myerc20.balanceOf(bidder1.address);
            const bidder2BlanceFinal = await myerc20.balanceOf(bidder2.address);
            const auctionBlanceFinal = await myerc20.balanceOf(auctionProxy.target);

            //验证bidder2是最高出价者
            const auctionInfoBid2 = await auctionProxy.getAuctionInfo();
            expect(auctionInfoBid2.highestBidder).to.equal(bidder2.address);
            expect(auctionInfoBid2.highestPrice).to.equal(secondBid);

            //验证bidder1退回后的余额，因为会扣除gas费，但gas费是用ETH支付的，
            //所以bidder1账户ETH前后的余额应该有些差距，但是竞拍的myerc20是完全返还的。
            expect(bidder1BlanceBefore).to.equal(bidder1BlanceFinal);
            // console.log("bidder1BlanceBefore:", bidder1BlanceBefore);
            // console.log("bidder1BlanceFinal:", bidder1BlanceFinal);
            // console.log("bidder2BlanceBefore:", bidder2BlanceBefore);
            // console.log("bidder2BlanceFinal:", bidder2BlanceFinal);
            // console.log("auctionBlanceBefore:", auctionBlanceBefore);
            // console.log("auctionBlanceFinal:", auctionBlanceFinal);
        });

        
    }); 
    
    
    describe("End Auction & Claim", function () { 
        const tokenId = 0;

        beforeEach(async () => {
            //卖家创建NFT，第一次tokenId=0
            await mynft.mint(seller.address, TOKEN_URI);
            //卖家授权
            await mynft.connect(seller).approve(auctionProxy, tokenId);

            //创建合约
            await auctionProxy.initialize(
                seller.address,
                mynft.target,
                tokenId,
                START_PRICE_ETH,
                10, //拍卖持续10秒
                ethers.ZeroAddress
            );

            //设置预言机
            await auctionProxy.setPriceFeeds(ethers.ZeroAddress, ethToUsdPriceFeed.target);
            await auctionProxy.setPriceFeeds(myerc20.target, erc20ToUsdPriceFeed.target);

        });

        
        it("Should end auction by anyone and NFT back to seller when no one bid", async function () { 
            //验证nft授权，此时拍卖开始，NFT授权给拍卖合约
            expect(await mynft.getApproved(tokenId)).to.equal(await auctionProxy.getAddress());
            
            //暂停10秒，等拍卖结束
            await new Promise(resolve => setTimeout(resolve, 10000)); 

            //任何人都可以在持续时间结束后，停止拍卖
            await auctionProxy.connect(bidder1).endAuction();

            //校验此时拍卖合约状态
            let auctionInfo = await auctionProxy.getAuctionInfo();
            expect(auctionInfo.isEnd).to.be.true;
            expect(auctionInfo.nftClaimed).to.be.false;
            expect(auctionInfo.paymentClaimed).to.be.false;

            //卖家收回授权
            await auctionProxy.claimNFTtoSeller();
            auctionInfo = await auctionProxy.getAuctionInfo();
            expect(auctionInfo.nftClaimed).to.be.true;

            //任何人都可以提取支付
            await auctionProxy.connect(seller).claimPayment();
            auctionInfo = await auctionProxy.getAuctionInfo();
            expect(auctionInfo.paymentClaimed).to.be.true;
            
        }); 

        
        it("Should allow winner to claim NFT and seller to receive payment", async function () {
            let auctionInfo = await auctionProxy.getAuctionInfo();
            const firstBid = ethers.parseEther("1.2");
            const secondBid = ethers.parseUnits("15", 18);
            
            //给竞拍者发放代币并授权
            await myerc20.transfer(bidder2.address, ethers.parseUnits("100", 18));
            await myerc20.connect(bidder2).approve(auctionProxy.target, ethers.parseUnits("100", 18));

            //const allowance = await myerc20.allowance(bidder2.address, auctionProxy.target);
            //console.log("Allowance:", ethers.formatUnits(allowance, 18)); // 应该是 100.0

            //账户初始余额
            const sellerInitBalEth = await ethers.provider.getBalance(seller.address);
            const bidder1InitBalEth = await ethers.provider.getBalance(bidder1.address);
            const bidder2InitBalEth = await ethers.provider.getBalance(bidder2.address);
            const auctionInitBalEth = await ethers.provider.getBalance(auctionProxy.target);
            const sellerInitBalErc20 = await myerc20.balanceOf(seller.address);
            const bidder1InitBalErc20 = await myerc20.balanceOf(bidder1.address);
            const bidder2InitBalErc20 = await myerc20.balanceOf(bidder2.address);
            const auctionInitBalErc20 = await myerc20.balanceOf(auctionProxy.target);

            //第一次竞拍
            await auctionProxy.connect(bidder1).bidByETH({value: firstBid});

            const sellerInitBalEth_after1 = await ethers.provider.getBalance(seller.address);
            const bidder1InitBalEth_after1 = await ethers.provider.getBalance(bidder1.address);
            const bidder2InitBalEth_after1 = await ethers.provider.getBalance(bidder2.address);
            const auctionInitBalEth_after1 = await ethers.provider.getBalance(auctionProxy.target);
            const sellerInitBalErc20_after1 = await myerc20.balanceOf(seller.address);
            const bidder1InitBalErc20_after1 = await myerc20.balanceOf(bidder1.address);
            const bidder2InitBalErc20_after1 = await myerc20.balanceOf(bidder2.address);
            const auctionInitBalErc20_after1 = await myerc20.balanceOf(auctionProxy.target);

            // console.log("sellerInitBalEth:", sellerInitBalEth);
            // console.log("sellerInitBalEth_after1:", sellerInitBalEth_after1);
            // console.log("bidder1InitBalEth:", bidder1InitBalEth);
            // console.log("bidder1InitBalEth_after1:", bidder1InitBalEth_after1);
            // console.log("bidder2InitBalEth:", bidder2InitBalEth);
            // console.log("bidder2InitBalEth_after1:", bidder2InitBalEth_after1);
            // console.log("auctionInitBalEth:", auctionInitBalEth);
            // console.log("auctionInitBalEth_after1:", auctionInitBalEth_after1);
            // console.log("sellerInitBalErc20:", sellerInitBalErc20);
            // console.log("sellerInitBalErc20_after1:", sellerInitBalErc20_after1);
            // console.log("bidder1InitBalErc20:", bidder1InitBalErc20);
            // console.log("bidder1InitBalErc20_after1:", bidder1InitBalErc20_after1);
            // console.log("bidder2InitBalErc20:", bidder2InitBalErc20);
            // console.log("bidder2InitBalErc20_after1:", bidder2InitBalErc20_after1);
            // console.log("auctionInitBalErc20:", auctionInitBalErc20);
            // console.log("auctionInitBalErc20_after1:", auctionInitBalErc20_after1);

            auctionInfo = await auctionProxy.getAuctionInfo();
            // console.log("bidder1.address:", bidder1.address);
            // console.log("auctionInfo.highestBidder:", auctionInfo.highestBidder);
            // console.log("auctionInfo.highestTokenAddress:", auctionInfo.highestTokenAddress);
            // console.log("auctionInfo.highestPrice:", auctionInfo.highestPrice);
            //console.log("auctionInfo_after1:", auctionInfo);

            //第二次竞拍
            await auctionProxy.connect(bidder2).bidByToken(myerc20.target, secondBid);

            const sellerInitBalEth_after2 = await ethers.provider.getBalance(seller.address);
            const bidder1InitBalEth_after2 = await ethers.provider.getBalance(bidder1.address);
            const bidder2InitBalEth_after2 = await ethers.provider.getBalance(bidder2.address);
            const auctionInitBalEth_after2 = await ethers.provider.getBalance(auctionProxy.target);
            const sellerInitBalErc20_after2 = await myerc20.balanceOf(seller.address);
            const bidder1InitBalErc20_after2 = await myerc20.balanceOf(bidder1.address);
            const bidder2InitBalErc20_after2 = await myerc20.balanceOf(bidder2.address);
            const auctionInitBalErc20_after2 = await myerc20.balanceOf(auctionProxy.target);

            // console.log("===第二次竞拍后===");
            // console.log("sellerInitBalEth_after1:", sellerInitBalEth_after1);
            // console.log("sellerInitBalEth_after2:", sellerInitBalEth_after2);
            // console.log("bidder1InitBalEth_after1:", bidder1InitBalEth_after1);
            // console.log("bidder1InitBalEth_after2:", bidder1InitBalEth_after2);
            // console.log("bidder2InitBalEth_after1:", bidder2InitBalEth_after1);
            // console.log("bidder2InitBalEth_after2:", bidder2InitBalEth_after2);
            // console.log("auctionInitBalEth_after1:", auctionInitBalEth_after1);
            // console.log("auctionInitBalEth_after2:", auctionInitBalEth_after2);
            // console.log("sellerInitBalErc20_after1:", sellerInitBalErc20_after1);
            // console.log("sellerInitBalErc20_after2:", sellerInitBalErc20_after2);
            // console.log("bidder1InitBalErc20_after1:", bidder1InitBalErc20_after1);
            // console.log("bidder1InitBalErc20_after2:", bidder1InitBalErc20_after2);
            // console.log("bidder2InitBalErc20_after1:", bidder2InitBalErc20_after1);
            // console.log("bidder2InitBalErc20_after2:", bidder2InitBalErc20_after2);
            // console.log("auctionInitBalErc20_after1:", auctionInitBalErc20_after1);
            // console.log("auctionInitBalErc20_after2:", auctionInitBalErc20_after2);

            auctionInfo = await auctionProxy.getAuctionInfo();
            // console.log("bidder1.address:", bidder2.address);
            // console.log("auctionInfo.highestBidder:", auctionInfo.highestBidder);
            // console.log("auctionInfo.highestTokenAddress:", auctionInfo.highestTokenAddress);
            // console.log("auctionInfo.highestPrice:", auctionInfo.highestPrice);
            //console.log("auctionInfo_after1:", auctionInfo);

            //暂停10秒，等拍卖结束
            await new Promise(resolve => setTimeout(resolve, 10000)); 
            //提前停止拍卖 //await auctionProxy.connect(seller).endAuctionEarly();
            await auctionProxy.connect(seller).endAuction();

            //出价最高者获得NFT
            expect(await mynft.ownerOf(tokenId)).to.equal(seller.address);
            await auctionProxy.connect(bidder2).claimNFT();
            auctionInfo = await auctionProxy.getAuctionInfo();
            expect(await mynft.ownerOf(tokenId)).to.equal(bidder2.address);
            expect(auctionInfo.nftClaimed).to.be.true;

            //卖家提取资金
            await auctionProxy.connect(seller).claimPayment();
            auctionInfo = await auctionProxy.getAuctionInfo();
            expect(auctionInfo.paymentClaimed).to.be.true;

            //竞拍结束后账户余额
            const sellerFinalBalEth = await ethers.provider.getBalance(seller.address);
            const bidder1FinalBalEth = await ethers.provider.getBalance(bidder1.address);
            const bidder2FinalBalEth = await ethers.provider.getBalance(bidder2.address);
            const auctionFinalBalEth = await ethers.provider.getBalance(auctionProxy.target);
            const sellerFinalBalErc20 = await myerc20.balanceOf(seller.address);
            const bidder1FinalBalErc20 = await myerc20.balanceOf(bidder1.address);
            const bidder2FinalBalErc20 = await myerc20.balanceOf(bidder2.address);
            const auctionFinalBalErc20 = await myerc20.balanceOf(auctionProxy.target);

            // console.log("===竞拍结束===");
            // console.log("sellerFinalBalEth:", sellerFinalBalEth);
            // console.log("bidder1FinalBalEth:", bidder1FinalBalEth);
            // console.log("bidder2FinalBalEth:", bidder2FinalBalEth);
            // console.log("sellerFinalBalErc20:", sellerFinalBalErc20);
            // console.log("bidder1FinalBalErc20:", bidder1FinalBalErc20);
            // console.log("bidder2FinalBalErc20:", bidder2FinalBalErc20);
            // console.log("auctionFinalBalEth:", auctionFinalBalEth);
            // console.log("auctionFinalBalErc20:", auctionFinalBalErc20);

            //seller收到的是ERC20代币
            expect(sellerFinalBalErc20).to.equal(sellerInitBalErc20 + secondBid);
            expect(sellerFinalBalEth).to.be.closeTo(sellerInitBalEth, ethers.parseEther("0.01"));
            //bidder1使用eth竞标，虽然没成功，但是会有gas费
            expect(bidder1InitBalEth).to.be.closeTo(bidder1FinalBalEth, ethers.parseEther("0.01"));
            expect(bidder1InitBalErc20).to.equal(bidder1FinalBalErc20);
            //bidder2使用erc20竞标，成功
            expect(bidder2FinalBalErc20).to.equal(bidder2InitBalErc20 - secondBid);
            expect(bidder2InitBalEth).to.be.closeTo(bidder2FinalBalEth, ethers.parseEther("0.01"));
        }); 


    }); 

});