// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract BinarySearch {

    function binarySearch(int[] memory array, int target) public pure returns (int) {

        int left = 0;
        int right = int(array.length - 1);

        while (left <= right) {
            int mid = left + (right - left) / 2; // 防止(left + right)溢出

            // 检查中间点是否是目标值
            if (array[uint(mid)] == target) {
                return mid; // 返回找到的索引
            }

            // 如果目标值大于中间点，则忽略左半边
            if (array[uint(mid)] < target) {
                left = mid + 1;
            } else {
                // 如果目标值小于中间点，则忽略右半边
                right = mid - 1;
            }

        }

        // 如果未找到目标值，返回-1
        return -1;

    }

}