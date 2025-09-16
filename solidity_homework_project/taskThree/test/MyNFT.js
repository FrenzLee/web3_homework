const { ethers } = require("hardhat");
const { expect } = require("chai");

describe("MyNFT Test", function () { 

    let MYNFT, mynft, owner, addr1, addr2;

    const MYNFT_NAME = "MyNFT";
    const MYNFT_SYMBOL = "MT";
    const TOKEN_URI_1 = "https://example.com/my-nft-1.json";
    const TOKEN_URI_2 = "https://example.com/my-nft-2.json";

    beforeEach(async () => {
        //获取账户地址
        [owner, addr1, addr2] = await ethers.getSigners();
        //获取合约工厂
        MYNFT = await ethers.getContractFactory("MyNFT");
        //部署合约
        mynft = await MYNFT.deploy(MYNFT_NAME, MYNFT_SYMBOL);
        //等待合约部署完成
        await mynft.waitForDeployment();
    });

    //测试用例分组
    describe("deployment", function () { 
        it("Should be deployed correctly", async function () { 
            expect(await mynft.name()).to.equal(MYNFT_NAME);
            expect(await mynft.symbol()).to.equal(MYNFT_SYMBOL);
            //mynft.owner() 智能合约本身 的“所有者”（通常是部署者或管理员）
            expect(await mynft.owner()).to.equal(owner.address);
        });
     });
    
    describe("minting", function () {
        it("Should mint a new NFT for addr1", async function () {
            //mintTx不是返回值tokenId，而是交易对象ContractTransactionResponse{...}
            const mintTx = await mynft.mint(addr1.address, TOKEN_URI_1)
            
            expect(await mynft.balanceOf(addr1.address)).to.equal(1);
            //tokenId = 0 的NFT应该归属于addr1
            expect(await mynft.ownerOf(0)).to.equal(addr1.address);
            expect(await mynft.tokenURI(0)).to.equal(TOKEN_URI_1);
        });

        it("Should return correct token ID", async function () {
            //静态调用合约函数，不会改变合约内部状态，只是模拟执行，所以返回值为tokenId，而不是交易对象
            const tokenId1 = await mynft.mint.staticCall(addr1.address, TOKEN_URI_1);
            expect(tokenId1).to.equal(0);

            await mynft.mint(addr1.address, TOKEN_URI_1);
            
            const tokenId2 = await mynft.mint.staticCall(addr1.address, TOKEN_URI_2);
            expect(tokenId2).to.equal(1);
        });

    });

    describe("transfer", function () {
        it("Should transfer NFT from addr1 to addr2 by addr1", async function () { 
            //addr1有2个NFT，tokenId分别是0，1，把 tokenId=1的NFT转给addr2
            await mynft.mint(addr1.address, TOKEN_URI_1);
            await mynft.mint(addr1.address, TOKEN_URI_2);
            await mynft.connect(addr1).transferFrom(addr1.address, addr2.address, 1);

            expect(await mynft.ownerOf(0)).to.equal(addr1.address);
            expect(await mynft.ownerOf(1)).to.equal(addr2.address);
            expect(await mynft.balanceOf(addr1.address)).to.equal(1);
            expect(await mynft.balanceOf(addr2.address)).to.equal(1);
            expect(await mynft.tokenURI(0)).to.equal(TOKEN_URI_1);
            expect(await mynft.tokenURI(1)).to.equal(TOKEN_URI_2);
        });

        it("Should transfer NFT from addr1 to addr2 by owner", async function () { 
            await mynft.mint(addr1.address, TOKEN_URI_1);
            //addr1将tokenId=0的NFT授权给owner
            await mynft.connect(addr1).approve(owner.address, 0);
            //owner将tokenId=1的NFT转给addr2
            await mynft.connect(owner).safeTransferFrom(addr1.address, addr2.address, 0);

            expect(await mynft.ownerOf(0)).to.equal(addr2.address);
            expect(await mynft.balanceOf(owner.address)).to.equal(0);
            expect(await mynft.balanceOf(addr1.address)).to.equal(0);
            expect(await mynft.balanceOf(addr2.address)).to.equal(1);
            expect(await mynft.tokenURI(0)).to.equal(TOKEN_URI_1);
        });
    });

    
});