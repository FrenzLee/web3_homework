// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

library SafeMath {
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        unchecked {
            uint256 c = a + b;
            if(c < a) revert("SafeMath: add overflow");
            return c;
        }
    }

    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a < b) revert("SafeMath: sub underflow");
        return a - b;
    }

    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) return 0;
        uint256 c = a * b;
        if (c / a != b) revert("SafeMath: mul overflow");
        return c;
    }

    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        if (b == 0) revert("SafeMath: div by zero");
        return a / b;
    }

}

contract OpMathCalV2 {
    using SafeMath for uint256;

    // 将多个变量打包到一个存储槽中
    // layout: [0-127: result][128-191: lastOperation][192-255: isInitialized]
    uint256 private _packedData;

    address private _owner;

    modifier onlyOwner() {
        require(msg.sender == _owner, "Only owner");
        _;
    }

    constructor() {
        _owner = msg.sender;
        _packedData = 1 << 255; // 初始化标志位
    }

    //打包
    function _packData(uint256 result, uint8 operation) private {
        _packedData = (result & ((1 << 128) - 1)) | // 结果放在0-128位
                     (uint256(operation) << 128) |  // 操作数放在第128-191位
                     (1 << 255); //1放在最高位，表示已经初始化了
    }

    //解包，获取结果
    function _unpackResult() private view returns (uint256) {
        return _packedData & ((1 << 128) - 1);
    }

    //解包，获取操作数
    function _unpackOperation() private view returns (uint8) {
        return uint8((_packedData >> 128) & 0xFF);
    }

    function add(uint256 a, uint256 b) external onlyOwner returns (uint256) {
        uint256 result = a.add(b);
        _packData(result, 0);
        return result;
    }

    function sub(uint256 a, uint256 b) external onlyOwner returns (uint256) {
        uint256 result = a.sub(b);
        _packData(result, 1);
        return result;
    }

    function mul(uint256 a, uint256 b) external onlyOwner returns (uint256) {
        uint256 result = a.mul(b);
        _packData(result, 2);
        return result;
    }

    function div(uint256 a, uint256 b) external onlyOwner returns (uint256) {
        uint256 result = a.div(b);
        _packData(result, 3);
        return result;
    }

    //平方
    function square(uint256 a) external onlyOwner returns (uint256) {
        uint256 result = a.mul(a);
        _packData(result, 4);
        return result;
    }

    //立方
    function cube(uint256 a) external onlyOwner returns (uint256) {
        uint256 result = a.mul(a).mul(a);
        _packData(result, 5);
        return result;
    }

    function getResult() external view returns (uint256) {
        return _unpackResult();
    }

    function getLastOperation() external view returns (uint8) {
        return _unpackOperation();
    }

    function isInitialized() external view returns (bool) {
        return (_packedData >> 255) == 1;
    }

    function reset() external onlyOwner {
        _packData(0, 0);
    }
}