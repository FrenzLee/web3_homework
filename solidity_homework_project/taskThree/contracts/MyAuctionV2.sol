// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "./MyAuction.sol";

contract MyAuctionV2 is MyAuction {
    function version() external pure returns (string memory) {
        return "v2.0.0";
    }
}

