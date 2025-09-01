// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract IntToRoma {
    //定义罗马值和整数，按照从小到大的顺序
    uint[] intArr = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1];
    string[] romaArr = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"];

    function intToRoma(uint inputNum) public view returns (string memory) {
        string memory outRoma;
        uint len = intArr.length;
        for (uint i = 0; i < len; i++) {
            while (inputNum >= intArr[i]) {
                inputNum -= intArr[i];
                outRoma = string(abi.encodePacked(outRoma, romaArr[i]));
            }
        }

        return outRoma;
    }

}
