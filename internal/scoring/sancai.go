package scoring

import "github.com/vogo/namer/internal/data"

// ScoreSanCai 三才五行配置评分（满分25）
func ScoreSanCai(wg WuGeResult) (float64, string) {
	tianWx := data.NumToWuXing(wg.TianGe)
	renWx := data.NumToWuXing(wg.RenGe)
	diWx := data.NumToWuXing(wg.DiGe)

	key := tianWx.String() + renWx.String() + diWx.String()
	jx, ok := data.SanCaiJiXiong[key]
	if !ok {
		return 12.5, key // 默认中间分
	}

	score := jiXiongToSanCaiScore(jx)
	return score / 100.0 * 25.0, key
}

func jiXiongToSanCaiScore(j data.JiXiong) float64 {
	switch j {
	case data.DaJi:
		return 100
	case data.ZhongJi:
		return 80
	case data.Ji:
		return 75
	case data.JiDuo:
		return 60
	case data.JiXiongBan:
		return 50
	case data.XiongDuo:
		return 30
	case data.DaXiong:
		return 10
	}
	return 50
}
