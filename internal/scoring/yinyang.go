package scoring

import "strings"

// ScoreYinYang 阴阳平衡评分（满分10）
// strokes: 姓名每个字的笔画数
func ScoreYinYang(strokes []int) (float64, string) {
	n := len(strokes)
	if n < 2 {
		return 5.0, "未知"
	}

	// 判断阴阳：奇数=阳, 偶数=阴
	pattern := make([]byte, n)
	for i, s := range strokes {
		if s%2 == 1 {
			pattern[i] = 'Y' // 阳
		} else {
			pattern[i] = 'I' // 阴
		}
	}
	var score float64
	if n == 2 {
		if pattern[0] != pattern[1] {
			score = 100 // 阴阳 或 阳阴
		} else {
			score = 40 // 纯阴或纯阳
		}
	} else {
		// 3字或以上
		allSame := true
		alternating := true
		for i := 1; i < n; i++ {
			if pattern[i] != pattern[0] {
				allSame = false
			}
			if pattern[i] == pattern[i-1] {
				alternating = false
			}
		}

		if alternating {
			score = 100 // 完全交替
		} else if allSame {
			score = 40 // 纯阴或纯阳
		} else {
			score = 80 // 有阴有阳但不完全交替
		}
	}

	// 生成阴阳描述
	var desc strings.Builder
	for _, p := range pattern {
		if p == 'Y' {
			desc.WriteString("阳")
		} else {
			desc.WriteString("阴")
		}
	}

	return score / 100.0 * 10.0, desc.String()
}
