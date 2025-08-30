package task_01

import "sort"

/*
56. 合并区间
以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间 。
*/

func Merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	// Sort intervals based on the starting point.
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	var merged [][]int

	for _, interval := range intervals {
		L, R := interval[0], interval[1]
		// If merged is empty or current interval does not overlap with the last one in merged,
		// simply append it to merged.
		if len(merged) == 0 || merged[len(merged)-1][1] < L {
			merged = append(merged, []int{L, R})
		} else {
			// Otherwise, there is an overlap, so we merge the current interval with the last one in merged.
			merged[len(merged)-1][1] = max(merged[len(merged)-1][1], R)
		}
	}

	return merged
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
