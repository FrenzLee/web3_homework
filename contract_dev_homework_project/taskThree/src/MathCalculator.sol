// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

contract MathCalculator {

    uint256 public result;
    address public owner;

    event CalPerformed(string operation, uint256 a, uint256 b, uint256 result);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can perform calculations");
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    
    function add(uint256 a, uint256 b) public onlyOwner returns (uint256) {
        result = a + b;
        emit CalPerformed("add", a, b, result);
        return result;
    }

    function sub(uint256 a, uint256 b) public onlyOwner returns (uint256) {
        require(a >= b, "Subtraction underflow");
        result = a - b;
        emit CalPerformed("sub", a, b, result);
        return result;
    }

    function mul(uint256 a, uint256 b) public onlyOwner returns (uint256) {
        result = a * b;
        emit CalPerformed("mul", a, b, result);
        return result;
    }

    function div(uint256 a, uint256 b) public onlyOwner returns (uint256) {
        require(b > 0, "Division by zero");
        result = a / b;
        emit CalPerformed("div", a, b, result);
        return result;
    }

    /**
     * 获取当前结果
     */
    function getResult() public view returns (uint256) {
        return result;
    }

    /**
     * 重置结果
     */
    function reset() public onlyOwner {
        result = 0;
    }

}