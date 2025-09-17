// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "./MyAuction.sol"; 
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/utils/Create2.sol";

contract MyAuctionFactory is Ownable {

    //用于生成预计算代理合约地址
    using Create2 for bytes32;

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

        //卖家创建拍卖时，不需要将NFT转给代理合约，只需要授权，需要先将NFT授权给预计算的代理合约
        //获取预计算代理合约
        //使用 CREATE2 计算唯一地址（salt 基于 nft + tokenId）
        bytes32 salt = keccak256(abi.encode(_nftAddress, _tokenId));
        //使用 CREATE2 部署合约
        auctionProxy = Create2.deploy(
            0, // 无附加 ETH
            salt,
            abi.encodePacked(type(ERC1967Proxy).creationCode, abi.encode(auctionImplementation, initData))
        );

        // 检查是否已授权
        address currentApproved = IERC721(_nftAddress).getApproved(_tokenId);
        require(currentApproved == auctionProxy, "NFT must be approved for the auction proxy");

        //记录信息
        auctionsMap[_nftAddress][_tokenId] = auctionProxy;
        allAuctions.push(auctionProxy);

        emit AuctionDeployed(msg.sender, auctionProxy, _nftAddress, _tokenId);
    }

    //预计算代理合约地址
    function predictAuctionAddress(
        address _seller,
        address _nftAddress, 
        uint256 _tokenId, 
        uint256 _startPrice,
        uint256 _duration,
        address _startTokenAddress
    ) public view returns (address) {
        bytes memory initData = abi.encodeWithSelector(
            MyAuction.initialize.selector,
            _seller,
            _nftAddress,
            _tokenId,
            _startPrice, 
            _duration, 
            _startTokenAddress
        );

        bytes memory initCode = abi.encodePacked(
            type(ERC1967Proxy).creationCode, 
            abi.encode(auctionImplementation, initData)
        );

        bytes32 salt = keccak256(abi.encode(_nftAddress, _tokenId));

        address pAuctionAddress = Create2.computeAddress(salt, keccak256(initCode));

        return pAuctionAddress;
    }


    //获取所有已创建的拍卖合约地址
    function getAllAuctions() external view returns (address[] memory) {
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