// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyERC20 is ERC20, Ownable {

    constructor(string memory name, string memory symbol, uint256 initialSupply) ERC20(name, symbol) Ownable(msg.sender) {
        _mint(msg.sender, initialSupply * 10**decimals());
    }

    function decimals() public pure virtual override returns (uint8) {
        return 18; // 或 6 等
    }

    /*
    function mint(address to, uint256 amount) external {
        _mint(to, amount* 10**decimals());
    } */

}