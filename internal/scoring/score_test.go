package scoring

import (
	"testing"

	"github.com/vogo/namer/internal/data"
)

func TestCalcScore(t *testing.T) {
	r := CalcScore("王", "明轩", 2024, 3, 15, 10, 30)

	// 总分在有效范围
	if r.Total < 0 || r.Total > 100 {
		t.Errorf("Total = %.1f, out of [0, 100]", r.Total)
	}
	// 各维度在有效范围
	if r.WuGeScore < 0 || r.WuGeScore > 30 {
		t.Errorf("WuGeScore = %.1f, out of [0, 30]", r.WuGeScore)
	}
	if r.SanCaiScore < 0 || r.SanCaiScore > 25 {
		t.Errorf("SanCaiScore = %.1f, out of [0, 25]", r.SanCaiScore)
	}
	if r.XiYongScore < 0 || r.XiYongScore > 20 {
		t.Errorf("XiYongScore = %.1f, out of [0, 20]", r.XiYongScore)
	}
	if r.WuXingScore < 0 || r.WuXingScore > 15 {
		t.Errorf("WuXingScore = %.1f, out of [0, 15]", r.WuXingScore)
	}
	if r.YinYangScore < 0 || r.YinYangScore > 10 {
		t.Errorf("YinYangScore = %.1f, out of [0, 10]", r.YinYangScore)
	}

	// 总分 = 各维度之和
	sum := r.WuGeScore + r.SanCaiScore + r.XiYongScore + r.WuXingScore + r.YinYangScore
	if diff := r.Total - sum; diff > 0.01 || diff < -0.01 {
		t.Errorf("Total %.1f != sum of parts %.1f", r.Total, sum)
	}

	// 笔画数
	if len(r.Strokes) != 3 {
		t.Errorf("Strokes length = %d, want 3", len(r.Strokes))
	}
	if r.Strokes[0] != 4 { // 王=4画
		t.Errorf("王 strokes = %d, want 4", r.Strokes[0])
	}

	// 五格
	if r.WuGe.TianGe != 5 {
		t.Errorf("TianGe = %d, want 5", r.WuGe.TianGe)
	}
	if r.WuGe.RenGe != 12 {
		t.Errorf("RenGe = %d, want 12", r.WuGe.RenGe)
	}

	// 三才描述非空
	if r.SanCaiDesc == "" {
		t.Error("SanCaiDesc should not be empty")
	}

	// 喜用神有效
	if r.XiYongShen == data.WuXingUnknown {
		t.Error("XiYongShen should not be unknown")
	}

	// 字五行
	if len(r.CharWuXing) != 3 {
		t.Errorf("CharWuXing length = %d, want 3", len(r.CharWuXing))
	}

	// 阴阳格局
	if r.YinYangPattern == "" {
		t.Error("YinYangPattern should not be empty")
	}
}

func TestCalcScoreSingleName(t *testing.T) {
	// 单字名
	r := CalcScore("张", "伟", 2024, 1, 1, 12, 0)
	if r.Total < 0 || r.Total > 100 {
		t.Errorf("Total = %.1f, out of [0, 100]", r.Total)
	}
	if len(r.Strokes) != 2 {
		t.Errorf("Strokes length = %d, want 2", len(r.Strokes))
	}
}

func TestCalcScoreMultipleNames(t *testing.T) {
	// 测试多种姓名组合都不panic且分数在有效范围
	names := []struct {
		last  string
		first string
	}{
		{"李", "浩然"},
		{"王", "明"},
		{"张", "伟"},
		{"赵", "子龙"},
		{"陈", "思"},
	}
	for _, n := range names {
		r := CalcScore(n.last, n.first, 2000, 6, 15, 8, 0)
		if r.Total < 0 || r.Total > 100 {
			t.Errorf("%s%s: Total = %.1f, out of [0, 100]", n.last, n.first, r.Total)
		}
	}
}

func TestCalcScoreDifferentBirthDates(t *testing.T) {
	// 不同出生日期影响喜用神，从而影响分数
	r1 := CalcScore("王", "明轩", 1990, 1, 1, 0, 0)
	r2 := CalcScore("王", "明轩", 2000, 6, 15, 12, 0)
	// 至少八字不同
	if r1.BaZi.String() == r2.BaZi.String() {
		t.Error("different birth dates should produce different BaZi")
	}
}

func TestPrintResult(t *testing.T) {
	// 只验证不 panic
	r := CalcScore("王", "明轩", 2024, 3, 15, 10, 30)
	PrintResult("王", "明轩", r)
}
