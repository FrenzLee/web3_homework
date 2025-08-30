package task_01

/*
136. 只出现一次的数字
给你一个 非空 整数数组 nums ，除了某个元素只出现一次以外，其余每个元素均出现两次。
找出那个只出现了一次的元素。
*/
func SingleNumber(nums []int) int {

	for i := range nums {
		a := nums[i]
		var count = 0
		for j := range nums {
			b := nums[j]
			if a == b {
				count += 1
			}
		}

		if count == 1 {
			return a
		}
	}

	return 0

}
