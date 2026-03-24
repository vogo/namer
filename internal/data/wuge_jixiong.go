package data

// WuGeJiXiong 1-81数理的吉凶
// 参考：81数理吉凶表
var WuGeJiXiong [82]JiXiong

func init() {
	// 大吉
	for _, n := range []int{1, 3, 5, 6, 11, 13, 15, 16, 21, 23, 24, 25, 29, 31, 32, 33, 35, 37, 39, 41, 45, 47, 48, 52, 57, 61, 63, 65, 67, 68, 73, 81} {
		WuGeJiXiong[n] = DaJi
	}
	// 吉
	for _, n := range []int{7, 8, 17, 18, 26, 38, 51, 55, 58, 71, 75} {
		WuGeJiXiong[n] = Ji
	}
	// 半吉 (JiDuo)
	for _, n := range []int{10, 16, 27, 30, 40, 42, 43, 50, 53, 57, 66, 72, 77, 78} {
		if WuGeJiXiong[n] == 0 && n != 0 {
			WuGeJiXiong[n] = JiDuo
		}
	}
	// 凶
	for _, n := range []int{2, 4, 9, 12, 14, 19, 20, 22, 26, 28, 34, 36, 44, 46, 49, 54, 56, 59, 60, 62, 64, 66, 69, 70, 72, 74, 76, 79, 80} {
		if WuGeJiXiong[n] == 0 && n != 0 {
			WuGeJiXiong[n] = DaXiong
		}
	}
}

// GetWuGeJiXiong 获取数理的吉凶，自动处理超过81的数理
func GetWuGeJiXiong(n int) JiXiong {
	if n <= 0 {
		return DaXiong
	}
	if n > 81 {
		n = n % 80
		if n == 0 {
			n = 80
		}
	}
	return WuGeJiXiong[n]
}
