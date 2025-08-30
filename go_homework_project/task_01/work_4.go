package task_01

/*
14. 最长公共前缀
编写一个函数来查找字符串数组中的最长公共前缀。
如果不存在公共前缀，返回空字符串 ""。
*/

// 横向扫描
func LongestCommonPrefix1(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	prefix := strs[0]  //假设第一个字符串为 最长公共前缀
	count := len(strs) //字符个数

	for i := 1; i < count; i++ {
		prefix = lcp(prefix, strs[i])
		if len(prefix) == 0 {
			break
		}
	}

	return prefix
}

func lcp(str1, str2 string) string {
	var min_leng int //最短字符串的长度
	if len(str1) < len(str2) {
		min_leng = len(str1)
	} else {
		min_leng = len(str2)
	}

	index := 0
	//比较两个字符串的每一个字符是否相同，取相同的部分
	for index < min_leng && str1[index] == str2[index] {
		index++
	}

	//截取两个字符串相同的部分
	return str1[:index]
}

// 2、纵向扫描
func LongestCommonPrefix2(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	//假设第一个字符串为公共前缀，循环其他字符串，对比相同位置的字符是否一致
	//i == len(strs[j]) 为了防止 strs[j][i] 报“数组越界”的错
	for i := 0; i < len(strs[0]); i++ {
		for j := 1; j < len(strs); j++ {
			if i == len(strs[j]) || string(strs[j][i]) != string(strs[0][i]) {
				return strs[0][:i]
			}
		}
	}

	return strs[0]
}
