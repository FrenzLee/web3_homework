// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import {ERC721URIStorage, ERC721} from "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyNFT is ERC721, ERC721URIStorage, Ownable {
    uint256 public _nextTokenId;

    constructor(string memory name, string memory symbol) ERC721(name, symbol) Ownable(msg.sender) {}

    //铸币
    function mint(address to, string memory tokenUri) public returns (uint256) {
        uint256 tokenId = _nextTokenId;
        _nextTokenId++;
        _mint(to, tokenId);
        _setTokenURI(tokenId, tokenUri);
        return tokenId;
    }  

    //获取tokenUri
    function tokenURI(uint256 tokenId) public view override(ERC721, ERC721URIStorage) returns (string memory) {
        return super.tokenURI(tokenId);
    }

    //接口是否支持
    function supportsInterface(bytes4 interfaceId) public view override(ERC721, ERC721URIStorage) returns (bool) {
        return super.supportsInterface(interfaceId);
    }


}