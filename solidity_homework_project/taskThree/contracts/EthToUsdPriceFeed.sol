// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

//模拟喂价合约
contract EthToUsdPriceFeed is AggregatorV3Interface, Ownable { 
    uint8 private decimals_;
    int256 private price_;
    uint80 private roundId_;

    constructor(uint8 _decimals, int256 _price) Ownable(msg.sender) {
        decimals_ = _decimals;
        price_ = _price;
        roundId_ = 1;
    }

    function decimals() external view override returns (uint8) {
        return decimals_;
    }

    function description() external pure override returns (string memory) {
        return "EthToUsdPriceFeed";
    }

    function version() external pure override returns (uint256) {
        return 1;
    }

    function getRoundData(
        uint80 _roundId
    ) external view returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound) {
        return (_roundId, price_, block.timestamp, block.timestamp, roundId_);
    }

    function latestRoundData() external view override returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    ) {
        return (roundId_, price_, block.timestamp, block.timestamp, roundId_);
    }

}