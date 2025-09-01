// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract RomaToInt {

    mapping (bytes1 => int) public romaToIntMap;

    //使用构造函数初始化
    constructor() {
        romaToIntMap['I'] = 1;
        romaToIntMap['V'] = 5;
        romaToIntMap['X'] = 10;
        romaToIntMap['L'] = 50;
        romaToIntMap['C'] = 100;
        romaToIntMap['D'] = 500;
        romaToIntMap['M'] = 1000;
    }

    function romaToInt(string memory roma) public view returns(int){
        int intSum = 0;
        bytes memory romaBytes = bytes(roma);
        uint len = romaBytes.length;

        for (uint i = 0; i < len; i++) {
            if (i<len-1 && romaToIntMap[romaBytes[i]]<romaToIntMap[romaBytes[i+1]]) {
                intSum -= romaToIntMap[romaBytes[i]];
            } else {
                intSum += romaToIntMap[romaBytes[i]];
            }
        }

        return intSum;

    }


}