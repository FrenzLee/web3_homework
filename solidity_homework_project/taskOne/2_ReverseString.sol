// SPDX-License-Identifier: MIT
pragma solidity ^0.8;


contract RevertseString {

    function revertString(string memory inputStr) public pure returns(string memory){

        bytes memory inputStrBytes = bytes(inputStr);
        uint strLen = inputStrBytes.length;
        bytes memory outputBytes = new bytes(strLen);

        for (uint i = 0; i < strLen; i++) {
            outputBytes[i] = inputStrBytes[strLen-1-i];
        }

        return string(outputBytes);

    }


}