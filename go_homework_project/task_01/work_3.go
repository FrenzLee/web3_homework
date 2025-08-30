package task_01

/*
20. 有效的括号
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
有效字符串需满足：
1、左括号必须用相同类型的右括号闭合。
2、左括号必须以正确的顺序闭合。
3、每个右括号都有一个对应的相同类型的左括号。
注意：s = "([)]"，"(]" -》false
s = "()[]{}"，"([])" -》true
*/
func IsValid(s string) bool {

	n := len(s)
	//如果是奇数，不是成对出现的，一定是false
	if n%2 == 1 {
		return false
	}

	//定义对应括号信息
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	//定义栈结构，用来存储左括号
	stack := []byte{}

	//开始处理
	for i := 0; i < n; i++ {
		if pairs[s[i]] > 0 {
			//如果是右括号，则判断栈顶元素是否对应
			//len(stack) == 0 表示第一个字符为右括号，无法按顺序匹配
			//stack[len(stack)-1] != pairs[s[i]] 表示栈顶元素与当前字符不是同类型括号，不匹配
			if len(stack) == 0 || stack[len(stack)-1] != pairs[s[i]] {
				return false
			}

			//移除栈顶元素，比如len(stack)=2，则新stack = stack[:1]，只有第一个元素，所以是移除了顶元素
			stack = stack[:len(stack)-1]
		} else {
			//如果是左括号，则入栈
			stack = append(stack, s[i])
		}
	}

	return len(stack) == 0

}
