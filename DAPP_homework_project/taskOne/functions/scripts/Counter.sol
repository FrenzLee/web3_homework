// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Counter {
    uint256 private count;

    event CountChanged(address indexed caller, uint256 indexed oldCount, uint256 newCount);

    constructor() {
        count = 0;
    }

    function getCurrentCount() public view returns (uint256) {
        return count;
    }

    function incrementOne() public {
        uint256 oldCount = count;
        count += 1;
        emit CountChanged(msg.sender, oldCount, count);
    }

    function resetZero() public {
        uint256 oldCount = count;
        count = 0;
        emit CountChanged(msg.sender, oldCount, count);
    }

}