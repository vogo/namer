package scoring

import (
	"testing"

	"github.com/vogo/namer/internal/data"
)

func TestXiYongRelationScore(t *testing.T) {
	// 木生火: Generates(木, 火) = true → 80
	if got := xiYongRelationScore(data.Mu, data.Huo); got != 80 {
		t.Errorf("xiYongRelationScore(木, 火) = %.0f, want 80", got)
	}
	// 火生土: Generates(火, 土) = true → 80 (名字五行生喜用神)
	if got := xiYongRelationScore(data.Huo, data.Tu); got != 80 {
		t.Errorf("xiYongRelationScore(火, 土) = %.0f, want 80", got)
	}
	// 喜用神生名字: Generates(金, 水) = true, so xiYong=金 charWx=水 → 60
	if got := xiYongRelationScore(data.Shui, data.Jin); got != 60 {
		t.Errorf("xiYongRelationScore(水, 金) = %.0f, want 60", got)
	}
	// 名字克喜用神: Overcomes(木, 土) = true → 20
	if got := xiYongRelationScore(data.Mu, data.Tu); got != 20 {
		t.Errorf("xiYongRelationScore(木, 土) = %.0f, want 20", got)
	}
	// 喜用神克名字: Overcomes(金, 木) = true → charWx=木 xiYong=金 → 10
	if got := xiYongRelationScore(data.Mu, data.Jin); got != 10 {
		t.Errorf("xiYongRelationScore(木, 金) = %.0f, want 10", got)
	}
	// 完全匹配
	if got := xiYongRelationScore(data.Jin, data.Jin); got != 100 {
		t.Errorf("xiYongRelationScore(金, 金) = %.0f, want 100", got)
	}
	// 未知
	if got := xiYongRelationScore(data.WuXingUnknown, data.Jin); got != 40 {
		t.Errorf("xiYongRelationScore(未知, 金) = %.0f, want 40", got)
	}
}

func TestScoreXiYong(t *testing.T) {
	// 空名字
	if got := ScoreXiYong(nil, data.Jin); got != 10.0 {
		t.Errorf("ScoreXiYong(nil) = %.1f, want 10.0", got)
	}
	// 未知喜用神
	if got := ScoreXiYong([]rune{'金'}, data.WuXingUnknown); got != 10.0 {
		t.Errorf("ScoreXiYong(unknown xiYong) = %.1f, want 10.0", got)
	}
	// 正常范围
	score := ScoreXiYong([]rune{'金', '银'}, data.Jin)
	if score < 0 || score > 20 {
		t.Errorf("ScoreXiYong = %.1f, out of [0, 20]", score)
	}
}

func TestInternalWuXingScore(t *testing.T) {
	// 五行相同
	if got := internalWuXingScore(data.Jin, data.Jin); got != 70 {
		t.Errorf("same element = %.0f, want 70", got)
	}
	// 前字生后字: 木生火
	if got := internalWuXingScore(data.Mu, data.Huo); got != 100 {
		t.Errorf("generates = %.0f, want 100", got)
	}
	// 后字生前字: 火生木? No. 水生木: b=水 a=木, Generates(水, 木) = true → 80
	if got := internalWuXingScore(data.Mu, data.Shui); got != 80 {
		t.Errorf("reverse generates = %.0f, want 80", got)
	}
	// 前字克后字: 木克土
	if got := internalWuXingScore(data.Mu, data.Tu); got != 20 {
		t.Errorf("overcomes = %.0f, want 20", got)
	}
	// 后字克前字: 金克木 → a=木 b=金, Overcomes(金, 木)=true → 30
	if got := internalWuXingScore(data.Mu, data.Jin); got != 30 {
		t.Errorf("reverse overcomes = %.0f, want 30", got)
	}
	// 未知
	if got := internalWuXingScore(data.WuXingUnknown, data.Jin); got != 50 {
		t.Errorf("unknown = %.0f, want 50", got)
	}
}

func TestScoreInternalWuXing(t *testing.T) {
	// 单字
	if got := ScoreInternalWuXing([]rune{'金'}); got != 7.5 {
		t.Errorf("single char = %.1f, want 7.5", got)
	}
	// 正常范围
	score := ScoreInternalWuXing([]rune{'水', '木', '火'}) // 水生木, 木生火 → 高分
	if score < 0 || score > 15 {
		t.Errorf("ScoreInternalWuXing = %.1f, out of [0, 15]", score)
	}
	// 连续相生应得高分
	if score < 13 {
		t.Errorf("水木火(连续相生) score = %.1f, want > 13", score)
	}
}
