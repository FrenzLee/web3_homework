package task_01

/*
9. 回文数
给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。
回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
例如，121 是回文，而 123 不是。
*/
func IsPalindrome(x int) bool {

	//负数肯定不是，最后一位是0的正整数也肯定不是
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}

	//反转数字
	revertNum := 0
	for x > revertNum {
		revertNum = revertNum*10 + x%10
		x /= 10
	}

	//偶数个数字相等，奇数个数字去掉最后一位后相等
	result := x == revertNum || x == revertNum/10

	return result

}
