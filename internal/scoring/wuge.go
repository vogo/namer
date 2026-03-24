package scoring

import "github.com/vogo/namer/internal/data"

// WuGeResult 五格计算结果
type WuGeResult struct {
	TianGe int // 天格
	RenGe  int // 人格
	DiGe   int // 地格
	ZongGe int // 总格
	WaiGe  int // 外格
}

// CalcWuGe 计算五格
// lastNameStrokes: 姓的各字笔画, firstNameStrokes: 名的各字笔画
func CalcWuGe(lastNameStrokes, firstNameStrokes []int) WuGeResult {
	singleLast := len(lastNameStrokes) == 1
	singleFirst := len(firstNameStrokes) == 1

	var r WuGeResult

	if singleLast && singleFirst {
		// 单姓单名
		ls := lastNameStrokes[0]
		fs := firstNameStrokes[0]
		r.TianGe = ls + 1
		r.RenGe = ls + fs
		r.DiGe = fs + 1
		r.ZongGe = ls + fs
		r.WaiGe = 2
	} else if singleLast && !singleFirst {
		// 单姓双名
		ls := lastNameStrokes[0]
		fs1 := firstNameStrokes[0]
		fs2 := firstNameStrokes[1]
		r.TianGe = ls + 1
		r.RenGe = ls + fs1
		r.DiGe = fs1 + fs2
		r.ZongGe = ls + fs1 + fs2
		r.WaiGe = r.ZongGe - r.RenGe + 1
	} else if !singleLast && singleFirst {
		// 复姓单名
		ls1 := lastNameStrokes[0]
		ls2 := lastNameStrokes[1]
		fs := firstNameStrokes[0]
		r.TianGe = ls1 + ls2
		r.RenGe = ls2 + fs
		r.DiGe = fs + 1
		r.ZongGe = ls1 + ls2 + fs
		r.WaiGe = r.ZongGe - r.RenGe + 1
	} else {
		// 复姓双名
		ls1 := lastNameStrokes[0]
		ls2 := lastNameStrokes[1]
		fs1 := firstNameStrokes[0]
		fs2 := firstNameStrokes[1]
		r.TianGe = ls1 + ls2
		r.RenGe = ls2 + fs1
		r.DiGe = fs1 + fs2
		r.ZongGe = ls1 + ls2 + fs1 + fs2
		r.WaiGe = r.ZongGe - r.RenGe
	}

	if r.WaiGe <= 0 {
		r.WaiGe = 2
	}

	return r
}

// ScoreWuGe 五格数理评分（满分30）
func ScoreWuGe(wg WuGeResult) float64 {
	tianScore := jiXiongToWuGeScore(data.GetWuGeJiXiong(wg.TianGe))
	renScore := jiXiongToWuGeScore(data.GetWuGeJiXiong(wg.RenGe))
	diScore := jiXiongToWuGeScore(data.GetWuGeJiXiong(wg.DiGe))
	zongScore := jiXiongToWuGeScore(data.GetWuGeJiXiong(wg.ZongGe))
	waiScore := jiXiongToWuGeScore(data.GetWuGeJiXiong(wg.WaiGe))

	weighted := renScore*0.35 + diScore*0.25 + zongScore*0.20 + waiScore*0.15 + tianScore*0.05
	return weighted / 100.0 * 30.0
}

func jiXiongToWuGeScore(j data.JiXiong) float64 {
	switch j {
	case data.DaJi:
		return 100
	case data.ZhongJi:
		return 90
	case data.Ji:
		return 85
	case data.JiDuo:
		return 65
	case data.JiXiongBan:
		return 50
	case data.XiongDuo:
		return 30
	case data.DaXiong:
		return 10
	}
	return 50
}
