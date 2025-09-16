// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "./MyAuction.sol"; 
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

contract MyAuctionFactory is Ownable {
    // 当前拍卖逻辑合约的实现地址
    address public auctionImplementation;
    // 实现版本号
    uint256 public implementationVersion;

    // NFT 地址 + TokenId => 拍卖合约地址
    mapping(address => mapping(uint256 => address)) public auctionsMap;
    // 所有创建过的拍卖合约地址列表
    address[] public allAuctions;

    // 事件：拍卖部署成功
    event AuctionDeployed(
        address indexed seller,
        address indexed auctionProxy,
        address nftAddress,
        uint256 tokenId
    );

    // 事件：拍卖逻辑合约升级
    event AuctionImplementationUpgraded(
        uint256 version,
        address oldImplementation,
        address newImplementation
    );

    //构造函数：传入所有者和初始拍卖逻辑实现地址
    constructor(
        address _owner,
        address _initialAuctionImplementation
    ) Ownable(_owner) {
        require(_owner != address(0), "Owner cannot be zero address");
        require(_initialAuctionImplementation != address(0), "Invalid implementation address");

        auctionImplementation = _initialAuctionImplementation;
        implementationVersion = 1;
    }

    //创建拍卖
    function createAuction(
        address _nftAddress,
        uint256 _tokenId,
        uint256 _startPrice,
        uint256 _duration,
        address _startTokenAddress
    ) external returns (address auctionProxy) { 
        require(_nftAddress != address(0), "Invalid NFT address");
        require(_startPrice > 0, "Start price must be greater than zero");
        require(_duration > 0, "Duration must be greater than zero");
        require(auctionsMap[_nftAddress][_tokenId] == address(0), "Auction already exists for this NFT");
        
        //卖家NFT转给合约工厂
        IERC721(_nftAddress).transferFrom(msg.sender, address(this), _tokenId);

        //构建 initialize 函数调用数据
        bytes memory initData = abi.encodeWithSelector(
            MyAuction.initialize.selector,
            msg.sender,
            _nftAddress,
            _tokenId,
            _startPrice,
            _duration,
            _startTokenAddress
        );

        //创建代理合约
        auctionProxy = address(new ERC1967Proxy(auctionImplementation, initData));

        //将NFT转给代理合约
        IERC721(_nftAddress).transferFrom(address(this), auctionProxy, _tokenId);

        //记录信息
        auctionsMap[_nftAddress][_tokenId] = auctionProxy;
        allAuctions.push(auctionProxy);

        emit AuctionDeployed(msg.sender, auctionProxy, _nftAddress, _tokenId);
    }

    //获取所有已创建的拍卖合约地址
    function getAuctions() external view returns (address[] memory) {
        return allAuctions;
    }

    //升级拍卖逻辑实现（仅所有者）
    function upgradeImplementation(address _newImplementation) external onlyOwner { 
        require(_newImplementation != address(0), "Invalid implementation address");
        require(_newImplementation != auctionImplementation, "Same as current implementation");

        address oldImplementation = auctionImplementation;
        auctionImplementation = _newImplementation;
        implementationVersion++;

        emit AuctionImplementationUpgraded(implementationVersion, oldImplementation, _newImplementation);
    }

    // 接收 ETH
    receive() external payable {}

}