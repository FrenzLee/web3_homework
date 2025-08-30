package task_01

import "fmt"

/*
26. 删除有序数组中的重复项
给你一个 非严格递增排列 的数组 nums ，请你 原地 删除重复出现的元素，使每个元素 只出现一次 ，
返回删除后数组的新长度。元素的 相对顺序 应该保持 一致 。然后返回 nums 中唯一元素的个数。

考虑 nums 的唯一元素的数量为 k ，你需要做以下事情确保你的题解可以被通过：
更改数组 nums ，使 nums 的前 k 个元素包含唯一元素，并按照它们最初在 nums 中出现的顺序排列。
nums 的其余元素与 nums 的大小不重要。
返回 k 。
*/

func RemoveDuplicates(nums []int) int {
	newNums := []int{}

	for i := 0; i < len(nums); i++ {

		if i == 0 {
			newNums = append(newNums, nums[i])
		}

		count := 0
		for j := 0; j < len(newNums); j++ {
			if nums[i] == newNums[j] {
				//newNums = append(newNums, nums[i])
				count++
			}
		}
		if count == 0 {
			newNums = append(newNums, nums[i])
		}

	}

	fmt.Println("原数组：", nums)
	fmt.Println("新数组：", newNums)

	return len(newNums)
}

func RemoveDuplicates1(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	slow := 1
	for fast := 1; fast < n; fast++ {
		if nums[fast] != nums[fast-1] {
			nums[slow] = nums[fast]
			slow++
		}
	}
	return slow
}
