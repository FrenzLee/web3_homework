// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

contract MyAuction is Initializable, UUPSUpgradeable, OwnableUpgradeable, ReentrancyGuardUpgradeable{
    
    //拍卖信息结构体
    struct AuctionInfo{
        address seller; //拍卖者
        address nftAddress;//拍卖品NFT地址
        uint256 tokenId;//拍卖品NFT的tokenId
        //uint256 auctionId; //拍卖ID
        uint256 startTime; //开始时间
        uint256 endTime; //结束时间
        uint256 startPrice; //起拍价
        uint256 highestPrice; //最高价
        address highestBidder; //最高价竞拍者
        bool isEnd; //是否结束
        bool nftClaimed; //拍卖结束后，NFT是否已经转给最高竞拍者
        bool paymentClaimed; //拍卖结束后，是否已经把钱转给拍卖者
        // 参与竞价的资产类型 0x 地址表示eth，其他地址表示erc20
        // 0x0000000000000000000000000000000000000000 表示eth
        address startTokenAddress;
        address highestTokenAddress;
    }

    using SafeERC20 for IERC20;
    //声明一个拍卖实体变量
    AuctionInfo public auctionInfo;
    //保存不同资产的喂价地址
    mapping(address => AggregatorV3Interface) public priceFeeds;
    // 支持的代币列表
    mapping(address => bool) public supportedTokens;

    //事件
    event AuctionCreated(
        address indexed seller,
        address indexed nftAddress,
        uint256 indexed tokenId,
        uint256 startPrice,
        uint256 endTime
    );

    event HighestBidIncreased(
        address indexed bidder,
        uint256 indexed tokenId,
        address highestTokenAddress,
        uint256 highestPrice
    );

    event AuctionEnded(
        address indexed seller,
        address indexed nftAddress,
        uint256 indexed tokenId,
        address highestBidder,
        uint256 highestPrice,
        address highestTokenAddress
    );

    event NFTClaimed(
        address indexed nftAddress,
        uint256 indexed tokenId,
        address highestBidder
    );

    event ReceivedETH(address sender, uint256 amount);

    //自定义修饰符
    modifier onlyBeforeEnd(){
        require(block.timestamp < auctionInfo.endTime && !auctionInfo.isEnd, "Auction has ended");
        _;
    } 

    modifier onlyAfterEnd(){
        require(block.timestamp >= auctionInfo.endTime, "Auction has not ended");
        _;
    }


    //确保禁用逻辑合约初始化功能，且代理合约只能通过initialize()函数初始化一次
    //使用UUPS不能有构造函数
    /*constructor() {
        _disableInitializers();
    }*/

    //初始化函数，只执行一次
    function initialize(
        address _seller, 
        address _nftAddress,//NFT所在的合约的地址
        uint256 _tokenId, 
        uint256 _startPrice,
        uint256 _duration,
        address _startTokenAddress //资产类型地址
    ) external initializer { 
        __Ownable_init(_seller);//OwnableUpgradeable初始化拍卖拥有人
        __ReentrancyGuard_init();//防重入攻击初始化
        __UUPSUpgradeable_init();//UUPS合约升级初始化

        //入参校验
        require(_seller != address(0), "Invalid seller address");
        require(_nftAddress != address(0), "Invalid NFT address");
        require(_tokenId >= 0, "Invalid token ID");
        require(_startPrice > 0, "Invalid start price");
        require(_duration > 0, "Invalid duration");
        //校验拍卖者的NFT是否授权给合约
        require(IERC721(_nftAddress).isApprovedForAll(_seller, address(this)) ||
                IERC721(_nftAddress).getApproved(_tokenId) == address(this),
                "NFT not approved for contract");

        auctionInfo = AuctionInfo({
            seller: _seller,
            nftAddress: _nftAddress,
            tokenId: _tokenId,
            startTime: block.timestamp,
            endTime: block.timestamp + _duration,
            startPrice: _startPrice,
            highestPrice: _startPrice,
            highestBidder: address(0),
            highestTokenAddress: address(0),
            isEnd: false,
            nftClaimed: false,
            paymentClaimed: false,
            startTokenAddress: _startTokenAddress
        });

        emit AuctionCreated(_seller, _nftAddress, _tokenId, _startPrice, auctionInfo.endTime);
    }

    //获取拍卖信息
    function getAuctionInfo() external view returns (AuctionInfo memory) {
        return auctionInfo;
    }

    //设置资产喂价地址 初始化后要设置
    //tokenAddress-资产地址，_priceFeed-资产对应喂价地址
    function setPriceFeeds(address tokenAddress, address _priceFeed) public {
        priceFeeds[tokenAddress] = AggregatorV3Interface(_priceFeed);
        supportedTokens[tokenAddress] = true;
    }

    //获取转换后的美元价，精度18，address(0) 表示 ETH
    function getPriceUSD(address tokenAddress, uint256 amount) public view returns (uint256) {
        require(supportedTokens[tokenAddress], "Token not supported");
        require(amount > 0, "Amount must be greater than zero");

        AggregatorV3Interface priceFeed = priceFeeds[tokenAddress];
        ( , int256 price, , , ) = priceFeed.latestRoundData();
        require(price > 0, "Invalid price");
        //喂价精度
        uint8 feedDecimals = priceFeed.decimals();
        
        //获取代币的精度
        uint8 tokenDecimals;
        if(tokenAddress == address(0)){
            tokenDecimals = 18;
        } else {
            tokenDecimals = IERC20Metadata(tokenAddress).decimals();
        }

        //将喂价价格转换为 18 位小数的 USD 价值（每单位代币）
        uint256 pricePerTokenIn18Decimals = uint256(price) * (10 ** (18 - feedDecimals));
        
        //计算总 USD 价值（18 位小数）
        //amount 是 tokenDecimals 位精度的整数,所以要除以 (10 ** tokenDecimals) 来对齐
        return (amount * pricePerTokenIn18Decimals) / (10 ** tokenDecimals);
    }

    //拍卖，使用ETH
    function bidByETH() external payable nonReentrant onlyBeforeEnd { 
        require(msg.value > 0, "Bid amount must be greater than zero");

        //竞拍价转换为美元
        uint256 bidUSD = getPriceUSD(address(0), msg.value);
        //起拍价转换为美元
        uint256 startPriceUSD = getPriceUSD(auctionInfo.startTokenAddress, auctionInfo.startPrice);
        //目前最高价转换为美元
        uint256 highestPriceUSD = getPriceUSD(auctionInfo.highestTokenAddress, auctionInfo.highestPrice);
        require(bidUSD >= startPriceUSD, "Bid amount must be higher than the start price");
        require(bidUSD > highestPriceUSD, "Bid amount must be higher than the current highest bid");

        //校验通过说明竞拍成功
        //退回之前的最高出价者，先退款，后更新状态
        _refundPreviousBidder();

        //更新拍卖信息
        auctionInfo.highestBidder = msg.sender;
        auctionInfo.highestPrice = msg.value;
        auctionInfo.highestTokenAddress = address(0);
        
        emit HighestBidIncreased(msg.sender, auctionInfo.tokenId, address(0), msg.value);
    } 


    //拍卖，使用ERC20代币
    function bidByToken(address payTokenAddress, uint256 amount) external nonReentrant onlyBeforeEnd { 
        require(payTokenAddress != address(0), "Invalid token address");
        require(supportedTokens[payTokenAddress], "Token not supported");
        require(amount > 0, "Bid amount must be greater than zero");

        //竞拍价转换为美元
        uint256 bidUSD = getPriceUSD(payTokenAddress, amount);
        //起拍价转换为美元
        uint256 startPriceUSD = getPriceUSD(auctionInfo.startTokenAddress, auctionInfo.startPrice);
        //目前最高价转换为美元
        uint256 highestPriceUSD = getPriceUSD(auctionInfo.highestTokenAddress, auctionInfo.highestPrice);
        require(bidUSD >= startPriceUSD, "Bid amount must be higher than the start price");
        require(bidUSD > highestPriceUSD, "Bid amount must be higher than the current highest bid");

        //将代币转移到拍卖合约
        IERC20(payTokenAddress).safeTransferFrom(msg.sender, address(this), amount);

        //校验通过说明竞拍成功
        //退回之前的最高出价者，先退款，后更新状态
        _refundPreviousBidder();

        //更新拍卖信息
        auctionInfo.highestBidder = msg.sender;
        auctionInfo.highestPrice = amount;
        auctionInfo.highestTokenAddress = payTokenAddress;
        
        emit HighestBidIncreased(msg.sender, auctionInfo.tokenId, payTokenAddress, amount);
    }

    //退回之前的最高出价者
    function _refundPreviousBidder() internal {
        if (auctionInfo.highestBidder != address(0)) {
            if (auctionInfo.highestTokenAddress == address(0)) {
                //退回ETH
                payable(auctionInfo.highestBidder).transfer(auctionInfo.highestPrice);
            } else {
                //退回ERC20代币
                IERC20(auctionInfo.highestTokenAddress).safeTransfer(auctionInfo.highestBidder, auctionInfo.highestPrice);
            }
        } 
    }

    //竞拍结束，此函数任何人都可以调用，只更新状态，onlyAfterEnd保证了拍卖的持续时间
    function endAuction() external onlyAfterEnd { 
        require(!auctionInfo.isEnd, "Auction has already ended");
        //修改状态
        auctionInfo.isEnd = true;

        //触发时间
        emit AuctionEnded(
            auctionInfo.seller, 
            auctionInfo.nftAddress, 
            auctionInfo.tokenId,
            auctionInfo.highestBidder,
            auctionInfo.highestPrice,
            auctionInfo.highestTokenAddress
        );

    }


    //拍卖结束后，NFT 转移给出价最高者，应该由出价最高者调用，
    //但是其他人也可调用，不过不管谁调用，NFT最终也只是转给highestBidder
    function claimNFT() external nonReentrant {
        require(auctionInfo.isEnd, "Auction: not ended");
        require(auctionInfo.highestBidder != address(0), "Auction: no winner");
        require(!auctionInfo.nftClaimed, "Auction: NFT already claimed");

        _transferNFT(auctionInfo.highestBidder);
    }

    //竞拍时间到达后，无人参与竞拍，退回NFT给卖家
    function claimNFTtoSeller() external nonReentrant {
        require(auctionInfo.isEnd, "Auction: not ended");
        require(auctionInfo.highestBidder == address(0), "Auction: no winner");
        require(!auctionInfo.nftClaimed, "Auction: NFT already claimed");

        //seller只是给了合约授权，只需要修改状态，不需要做其他事情
        auctionInfo.nftClaimed = true;

        //收回NFT对拍卖合约的授权
        //IERC721(auctionInfo.nftAddress).approve(address(0), auctionInfo.tokenId);
        
        //_transferNFT(auctionInfo.seller);
    }

    //转移NFT
    function _transferNFT(address to) internal { 
        require(!auctionInfo.nftClaimed, "NFT has already been claimed");
        //从seller转出NFT到to
        IERC721(auctionInfo.nftAddress).safeTransferFrom(auctionInfo.seller, to, auctionInfo.tokenId);
        //转换状态
        auctionInfo.nftClaimed = true;
        //触发事件
        emit NFTClaimed(auctionInfo.nftAddress, auctionInfo.tokenId, to);
    }

    //转移资金
    function claimPayment() external nonReentrant {
        require(auctionInfo.isEnd, "Auction: not ended");
        require(!auctionInfo.paymentClaimed, "Auction: payment already claimed");

        //竞拍时间到达后，无人参与竞拍，就不能转账，因为seller并没有向拍卖合约转过资产
        if (auctionInfo.highestBidder != address(0)) {
            uint256 sellAmount = auctionInfo.highestPrice;
            address payTokenAddress = auctionInfo.highestTokenAddress;

            if (payTokenAddress == address(0)) {
                //将ETH转账给卖家
                payable(auctionInfo.seller).transfer(sellAmount);
            } else {
                //将代币转账给卖家
                IERC20(payTokenAddress).safeTransfer(auctionInfo.seller, sellAmount);
            }
        }

        auctionInfo.paymentClaimed = true;
    }

    //授权升级
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    //可以接受ETH
    receive() external payable {
        emit ReceivedETH(msg.sender, msg.value);
    }

    // 管理员提取意外收到的 ETH
    function withdrawStuckETH(address to, uint256 amount) external onlyOwner {
        require(address(this).balance >= amount, "Insufficient balance");
        payable(to).transfer(amount);
    }

}