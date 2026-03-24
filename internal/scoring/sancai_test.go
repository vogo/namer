package scoring

import (
	"testing"

	"github.com/vogo/namer/internal/data"
)

func TestScoreSanCai(t *testing.T) {
	tests := []struct {
		name    string
		wg      WuGeResult
		wantMin float64
		wantMax float64
	}{
		{
			name:    "木木木=大吉",
			wg:      WuGeResult{TianGe: 11, RenGe: 11, DiGe: 11}, // 1,2→木
			wantMin: 24, wantMax: 25,
		},
		{
			name:    "水火火=大凶",
			wg:      WuGeResult{TianGe: 9, RenGe: 3, DiGe: 3}, // 9→水, 3→火
			wantMin: 2, wantMax: 3,
		},
	}
	for _, tt := range tests {
		score, desc := ScoreSanCai(tt.wg)
		if score < tt.wantMin || score > tt.wantMax {
			t.Errorf("%s: ScoreSanCai = %.1f (desc=%s), want [%.0f, %.0f]",
				tt.name, score, desc, tt.wantMin, tt.wantMax)
		}
		if desc == "" {
			t.Errorf("%s: desc should not be empty", tt.name)
		}
	}
}

func TestScoreSanCaiRange(t *testing.T) {
	// 对各种五格组合验证分数在 [0, 25] 范围内
	geValues := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 20}
	for _, tg := range geValues {
		for _, rg := range geValues {
			for _, dg := range geValues {
				wg := WuGeResult{TianGe: tg, RenGe: rg, DiGe: dg}
				score, _ := ScoreSanCai(wg)
				if score < 0 || score > 25 {
					t.Errorf("ScoreSanCai(%d,%d,%d) = %.1f, out of [0,25]", tg, rg, dg, score)
				}
			}
		}
	}
}

func TestJiXiongToSanCaiScore(t *testing.T) {
	tests := []struct {
		jx   data.JiXiong
		want float64
	}{
		{data.DaJi, 100},
		{data.ZhongJi, 80},
		{data.Ji, 75},
		{data.JiDuo, 60},
		{data.JiXiongBan, 50},
		{data.XiongDuo, 30},
		{data.DaXiong, 10},
		{data.JiXiong(99), 50}, // default
	}
	for _, tt := range tests {
		if got := jiXiongToSanCaiScore(tt.jx); got != tt.want {
			t.Errorf("jiXiongToSanCaiScore(%v) = %.0f, want %.0f", tt.jx, got, tt.want)
		}
	}
}
