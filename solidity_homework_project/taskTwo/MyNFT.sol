// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

import {ERC721URIStorage, ERC721} from "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/*
作业2：在测试网上发行一个图文并茂的 NFT
任务目标 
1. 使用 Solidity 编写一个符合 ERC721 标准的 NFT 合约。
2. 将图文数据上传到 IPFS，生成元数据链接。
3. 将合约部署到以太坊测试网（如 Goerli 或 Sepolia）。
4. 铸造 NFT 并在测试网环境中查看。
任务步骤 
1. 编写 NFT 合约
  ○ 使用 OpenZeppelin 的 ERC721 库编写一个 NFT 合约。
  ○ 合约应包含以下功能：
  ○ 构造函数：设置 NFT 的名称和符号。
  ○ mintNFT 函数：允许用户铸造 NFT，并关联元数据链接（tokenURI）。
  ○ 在 Remix IDE 中编译合约。
2. 准备图文数据
  ○ 准备一张图片，并将其上传到 IPFS（可以使用 Pinata 或其他工具）。
  ○ 创建一个 JSON 文件，描述 NFT 的属性（如名称、描述、图片链接等）。
  ○ 将 JSON 文件上传到 IPFS，获取元数据链接。
  ○ JSON文件参考 https://docs.opensea.io/docs/metadata-standards
3. 部署合约到测试网
  ○ 在 Remix IDE 中连接 MetaMask，并确保 MetaMask 连接到 Goerli 或 Sepolia 测试网。
  ○ 部署 NFT 合约到测试网，并记录合约地址。
4. 铸造 NFT
  ○ 使用 mintNFT 函数铸造 NFT：
  ○ 在 recipient 字段中输入你的钱包地址。
  ○ 在 tokenURI 字段中输入元数据的 IPFS 链接。
  ○ 在 MetaMask 中确认交易。
5. 查看 NFT
  ○ 打开 OpenSea 测试网 或 Etherscan 测试网。
  ○ 连接你的钱包，查看你铸造的 NFT。

元数据地址：https://ipfs.io/ipfs/bafybeifeph7vunjrmwioeqgiuboreolvuzm2sbw2vdupiexvplwkv4v3dy/MyNFT1.json
图片地址：https://ipfs.io/ipfs/bafybeiejsaitzilbhhcwe7hmu6ud7vcah4hjf3dhluozkjfafnhjwntvnq/MyNFT1.JPG
合约创建测试区块链地址：https://sepolia.etherscan.io/tx/0xa12ba61284a90b0b3aed40095cab6bbea7a4c6a3767712b9976d1c8db836e5e2
铸币测试区块链地址：https://sepolia.etherscan.io/tx/0x3b1653f267a1c112991a812fd7aac3e8119cfcf049f25966664be93965845cb9
*/


contract MyNFT2 is ERC721URIStorage, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;

    constructor() ERC721("LlNFT2", "LNP2") Ownable(msg.sender) {}

    //只有所有者可以铸造
    function mintNFT(address recipient, string memory tokenURI) public onlyOwner returns (uint256) {
        _tokenIds.increment();//id自增
        uint256 newTokenId = _tokenIds.current();//获取自增后的新id

        _mint(recipient, newTokenId);//铸造
        _setTokenURI(newTokenId, tokenURI);//关联图片

        return newTokenId;
    }
}
