package scoring

import "github.com/vogo/namer/internal/data"

// ScoreXiYong 名字五行与喜用神匹配评分（满分20）
func ScoreXiYong(firstNameChars []rune, xiYong data.WuXing) float64 {
	if xiYong == data.WuXingUnknown || len(firstNameChars) == 0 {
		return 10.0 // 默认中间分
	}

	total := 0.0
	for _, c := range firstNameChars {
		wx := data.CharWuXing(c)
		total += xiYongRelationScore(wx, xiYong)
	}
	avg := total / float64(len(firstNameChars))
	return avg / 100.0 * 20.0
}

func xiYongRelationScore(charWx, xiYong data.WuXing) float64 {
	if charWx == data.WuXingUnknown {
		return 40 // 未知五行，给中间分
	}
	if charWx == xiYong {
		return 100 // 完全匹配
	}
	if data.Generates(charWx, xiYong) {
		return 80 // 名字五行生喜用神
	}
	if data.Generates(xiYong, charWx) {
		return 60 // 喜用神生名字五行
	}
	if data.Overcomes(charWx, xiYong) {
		return 20 // 名字五行克喜用神
	}
	if data.Overcomes(xiYong, charWx) {
		return 10 // 喜用神克名字五行
	}
	return 40 // 无直接关系
}

// ScoreInternalWuXing 名字内部五行生克评分（满分15）
func ScoreInternalWuXing(allChars []rune) float64 {
	if len(allChars) < 2 {
		return 7.5 // 默认中间分
	}

	total := 0.0
	pairs := 0
	for i := 0; i < len(allChars)-1; i++ {
		wx1 := data.CharWuXing(allChars[i])
		wx2 := data.CharWuXing(allChars[i+1])
		total += internalWuXingScore(wx1, wx2)
		pairs++
	}

	avg := total / float64(pairs)
	return avg / 100.0 * 15.0
}

func internalWuXingScore(a, b data.WuXing) float64 {
	if a == data.WuXingUnknown || b == data.WuXingUnknown {
		return 50
	}
	if a == b {
		return 70 // 五行相同
	}
	if data.Generates(a, b) {
		return 100 // 前字生后字
	}
	if data.Generates(b, a) {
		return 80 // 后字生前字
	}
	if data.Overcomes(a, b) {
		return 20 // 前字克后字
	}
	if data.Overcomes(b, a) {
		return 30 // 后字克前字
	}
	return 50
}
