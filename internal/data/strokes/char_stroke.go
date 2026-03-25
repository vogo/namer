package strokes

// CharStroke 查询汉字的康熙笔画数，未找到返回 0
func CharStroke(c rune) int {
	if Strokes1[c] {
		return 1
	}
	if Strokes2[c] {
		return 2
	}
	if Strokes3[c] {
		return 3
	}
	if Strokes4[c] {
		return 4
	}
	if Strokes5[c] {
		return 5
	}
	if Strokes6[c] {
		return 6
	}
	if Strokes7[c] {
		return 7
	}
	if Strokes8[c] {
		return 8
	}
	if Strokes9[c] {
		return 9
	}
	if Strokes10[c] {
		return 10
	}
	if Strokes11[c] {
		return 11
	}
	if Strokes12[c] {
		return 12
	}
	if Strokes13[c] {
		return 13
	}
	if Strokes14[c] {
		return 14
	}
	if Strokes15[c] {
		return 15
	}
	if Strokes16[c] {
		return 16
	}
	if Strokes17[c] {
		return 17
	}
	if Strokes18[c] {
		return 18
	}
	if Strokes19[c] {
		return 19
	}
	if Strokes20[c] {
		return 20
	}
	if Strokes21[c] {
		return 21
	}
	if Strokes22[c] {
		return 22
	}
	if Strokes23[c] {
		return 23
	}
	if Strokes24[c] {
		return 24
	}
	if Strokes25[c] {
		return 25
	}
	if Strokes26[c] {
		return 26
	}
	if Strokes27[c] {
		return 27
	}
	if Strokes28[c] {
		return 28
	}
	if Strokes29[c] {
		return 29
	}
	if Strokes30[c] {
		return 30
	}
	if Strokes32[c] {
		return 32
	}
	return 0
}
