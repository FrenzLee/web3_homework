// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract MergeSortedArray{

    function MergeArr(uint[] memory arr1, uint[] memory arr2) public pure returns (uint[] memory) {
        uint len1 = arr1.length;
        uint len2 = arr2.length;
        uint[] memory outPutArr = new uint[](len1+len2);

        uint i = 0;
        uint j = 0;
        uint k = 0;

        while (i < len1 && j < len2) {
            if (arr1[i] <= arr2[j]) {
                outPutArr[k] = arr1[i];
                i++;
            } else {
                outPutArr[k] = arr2[j];
                j++;
            }
            k++;
        }

        while (i < arr1.length) {
            outPutArr[k] = arr1[i];
            i++;
            k++;
        }

        while (j < arr2.length) {
            outPutArr[k] = arr2[j];
            j++;
            k++;
        }

        return outPutArr;

    }

}