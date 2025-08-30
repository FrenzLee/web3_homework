package task_01

/*
66. 加一
给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。
这些数字按从左到右，从最高位到最低位排列。这个大整数不包含任何前导 0。
将大整数加 1，并返回结果的数字数组。
*/

func PlusOne(digits []int) []int {
	var number int
	for i := 0; i < len(digits); i++ {
		number = number*10 + digits[i]
	}

	number++

	var numRes []int
	for number > 0 {
		numRes = append(numRes, number%10)
		number /= 10
	}

	result := make([]int, len(numRes))
	for i := 0; i < len(numRes); i++ { //0,1,2
		result[i] = numRes[len(numRes)-i-1]
	}

	return result
}

func PlusOne1(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		if digits[i] != 9 {
			digits[i]++

			for j := i + 1; j < n; j++ {
				digits[j] = 0
			}

			return digits
		}
	}
	// digits 中所有的元素均为 9

	digits = make([]int, n+1)
	digits[0] = 1
	return digits
}
