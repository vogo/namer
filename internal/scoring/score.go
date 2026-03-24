package scoring

import (
	"fmt"
	"strings"

	"github.com/vogo/namer/internal/data"
)

// ScoreResult 评分结果
type ScoreResult struct {
	Total        float64
	WuGeScore    float64
	SanCaiScore  float64
	XiYongScore  float64
	WuXingScore  float64
	YinYangScore float64

	// 详细信息
	Strokes        []int
	WuGe           WuGeResult
	SanCaiDesc     string
	BaZi           BaZiResult
	XiYongShen     data.WuXing
	CharWuXing     []data.WuXing
	YinYangPattern string
}

// CalcScore 综合评分
func CalcScore(lastName, firstName string, year, month, day, hour, minute int) ScoreResult {
	lastRunes := []rune(lastName)
	firstRunes := []rune(firstName)
	allRunes := append(lastRunes, firstRunes...)

	// 1. 计算笔画
	strokes := make([]int, len(allRunes))
	for i, c := range allRunes {
		strokes[i] = data.CharStroke(c)
	}

	lastStrokes := strokes[:len(lastRunes)]
	firstStrokes := strokes[len(lastRunes):]

	// 2. 五格
	wg := CalcWuGe(lastStrokes, firstStrokes)

	// 3. 五格评分
	wuGeScore := ScoreWuGe(wg)

	// 4. 三才评分
	sanCaiScore, sanCaiDesc := ScoreSanCai(wg)

	// 5. 八字 & 喜用神
	bz := CalcBaZi(year, month, day, hour)
	xiYong := CalcXiYongShen(bz)

	// 6. 喜用神匹配评分
	xiYongScore := ScoreXiYong(firstRunes, xiYong)

	// 7. 内部五行评分
	charWuXing := make([]data.WuXing, len(allRunes))
	for i, c := range allRunes {
		charWuXing[i] = data.CharWuXing(c)
	}
	wuXingScore := ScoreInternalWuXing(allRunes)

	// 8. 阴阳评分
	yinYangScore, yinYangPattern := ScoreYinYang(strokes)

	total := wuGeScore + sanCaiScore + xiYongScore + wuXingScore + yinYangScore

	return ScoreResult{
		Total:          total,
		WuGeScore:      wuGeScore,
		SanCaiScore:    sanCaiScore,
		XiYongScore:    xiYongScore,
		WuXingScore:    wuXingScore,
		YinYangScore:   yinYangScore,
		Strokes:        strokes,
		WuGe:           wg,
		SanCaiDesc:     sanCaiDesc,
		BaZi:           bz,
		XiYongShen:     xiYong,
		CharWuXing:     charWuXing,
		YinYangPattern: yinYangPattern,
	}
}

// PrintResult 打印评分结果
func PrintResult(lastName, firstName string, r ScoreResult) {
	allRunes := []rune(lastName + firstName)

	fmt.Printf("\n姓名: %s%s\n", lastName, firstName)
	fmt.Printf("总分: %.1f / 100\n\n", r.Total)

	fmt.Println("┌──────────────┬──────┬──────┐")
	fmt.Println("│ 评分维度     │ 得分 │ 满分 │")
	fmt.Println("├──────────────┼──────┼──────┤")
	fmt.Printf("│ 五格数理     │ %4.1f │  30  │\n", r.WuGeScore)
	fmt.Printf("│ 三才配置     │ %4.1f │  25  │\n", r.SanCaiScore)
	fmt.Printf("│ 喜用神匹配   │ %4.1f │  20  │\n", r.XiYongScore)
	fmt.Printf("│ 内部五行     │ %4.1f │  15  │\n", r.WuXingScore)
	fmt.Printf("│ 阴阳平衡     │ %4.1f │  10  │\n", r.YinYangScore)
	fmt.Println("└──────────────┴──────┴──────┘")

	// 笔画
	fmt.Printf("\n康熙笔画: ")
	for i, c := range allRunes {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%c(%d)", c, r.Strokes[i])
	}
	fmt.Println()

	// 五格
	fmt.Printf("五格: 天格%d 人格%d 地格%d 总格%d 外格%d\n",
		r.WuGe.TianGe, r.WuGe.RenGe, r.WuGe.DiGe, r.WuGe.ZongGe, r.WuGe.WaiGe)

	// 三才
	fmt.Printf("三才: %s → %s\n", r.SanCaiDesc,
		data.SanCaiJiXiong[r.SanCaiDesc].String())

	// 八字
	fmt.Printf("八字: %s\n", r.BaZi.String())
	fmt.Printf("喜用神: %s\n", r.XiYongShen.String())

	// 字五行
	wxStrs := make([]string, len(r.CharWuXing))
	for i, wx := range r.CharWuXing {
		wxStrs[i] = wx.String()
	}
	fmt.Printf("字五行: %s\n", strings.Join(wxStrs, " "))

	// 阴阳
	fmt.Printf("阴阳: %s\n", r.YinYangPattern)
}
