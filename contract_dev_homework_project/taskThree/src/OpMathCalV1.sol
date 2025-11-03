// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

contract OpMathCalV1 {

    uint256 private _result;
    address private _owner;

    struct CalData {
        uint128 operaA;
        uint128 operaB;
    }
    CalData private _calDate;

    modifier onlyOwner() {
        require(msg.sender == _owner, "Only owner");
        _;
    }

    constructor() {
        _owner = msg.sender;
    }

    function add(uint256 a, uint256 b) external onlyOwner returns (uint256) {
        assembly {
            if lt(add(a, b), a) { revert(0, 0) } // 检查溢出
        }
        _result = a + b;
        _calDate = CalData(uint128(a), uint128(b));
        return _result;
    }

    function sub(uint256 a, uint256 b) external onlyOwner returns (uint256) {
        require(a >= b, "Underflow");
        _result = a - b;
        _calDate = CalData(uint128(a), uint128(b));
        return _result;
    }

    function mul(uint256 a, uint256 b) external onlyOwner returns (uint256) {
        if (a == 0) return 0;
        require(b == 0 || a <= type(uint256).max / b, "Overflow");
        _result = a * b;
        _calDate = CalData(uint128(a), uint128(b));
        return _result;
    }

    function div(uint256 a, uint256 b) external onlyOwner returns (uint256) {
        require(b > 0, "Division by zero");
        _result = a / b;
        _calDate = CalData(uint128(a), uint128(b));
        return _result;
    }


    function getResult() external view returns (uint256) {
        return _result;
    }

    function reset() external onlyOwner {
        _result = 0;
        delete _calDate;
    }

    function getLastCalculation() external view returns (uint128 a, uint128 b) {
        a = _calDate.operaA;
        b = _calDate.operaB;
    }



}