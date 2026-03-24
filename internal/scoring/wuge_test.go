package scoring

import (
	"testing"

	"github.com/vogo/namer/internal/data"
)

func TestCalcWuGeSingleLastSingleFirst(t *testing.T) {
	// 单姓单名：刘(15) 江(7)
	wg := CalcWuGe([]int{15}, []int{7})
	if wg.TianGe != 16 {
		t.Errorf("TianGe = %d, want 16", wg.TianGe)
	}
	if wg.RenGe != 22 {
		t.Errorf("RenGe = %d, want 22", wg.RenGe)
	}
	if wg.DiGe != 8 {
		t.Errorf("DiGe = %d, want 8", wg.DiGe)
	}
	if wg.ZongGe != 22 {
		t.Errorf("ZongGe = %d, want 22", wg.ZongGe)
	}
	if wg.WaiGe != 2 {
		t.Errorf("WaiGe = %d, want 2", wg.WaiGe)
	}
}

func TestCalcWuGeSingleLastDoubleFirst(t *testing.T) {
	// 单姓双名：王(4) 明(8) 轩(10)
	wg := CalcWuGe([]int{4}, []int{8, 10})
	if wg.TianGe != 5 {
		t.Errorf("TianGe = %d, want 5", wg.TianGe)
	}
	if wg.RenGe != 12 {
		t.Errorf("RenGe = %d, want 12", wg.RenGe)
	}
	if wg.DiGe != 18 {
		t.Errorf("DiGe = %d, want 18", wg.DiGe)
	}
	if wg.ZongGe != 22 {
		t.Errorf("ZongGe = %d, want 22", wg.ZongGe)
	}
	if wg.WaiGe != 11 {
		t.Errorf("WaiGe = %d, want 11", wg.WaiGe)
	}
}

func TestCalcWuGeDoubleLastSingleFirst(t *testing.T) {
	// 复姓单名：司马(5+10) 光(6)
	wg := CalcWuGe([]int{5, 10}, []int{6})
	if wg.TianGe != 15 {
		t.Errorf("TianGe = %d, want 15", wg.TianGe)
	}
	if wg.RenGe != 16 {
		t.Errorf("RenGe = %d, want 16", wg.RenGe)
	}
	if wg.DiGe != 7 {
		t.Errorf("DiGe = %d, want 7", wg.DiGe)
	}
	if wg.ZongGe != 21 {
		t.Errorf("ZongGe = %d, want 21", wg.ZongGe)
	}
	if wg.WaiGe != 6 {
		t.Errorf("WaiGe = %d, want 6", wg.WaiGe)
	}
}

func TestCalcWuGeDoubleLastDoubleFirst(t *testing.T) {
	// 复姓双名
	wg := CalcWuGe([]int{5, 10}, []int{22, 15})
	if wg.TianGe != 15 {
		t.Errorf("TianGe = %d, want 15", wg.TianGe)
	}
	if wg.RenGe != 32 {
		t.Errorf("RenGe = %d, want 32", wg.RenGe)
	}
	if wg.DiGe != 37 {
		t.Errorf("DiGe = %d, want 37", wg.DiGe)
	}
	if wg.ZongGe != 52 {
		t.Errorf("ZongGe = %d, want 52", wg.ZongGe)
	}
	if wg.WaiGe != 20 {
		t.Errorf("WaiGe = %d, want 20", wg.WaiGe)
	}
}

func TestCalcWuGeWaiGeMinimum(t *testing.T) {
	wg := CalcWuGe([]int{1, 1}, []int{1, 1})
	if wg.WaiGe < 2 {
		t.Errorf("WaiGe should be at least 2, got %d", wg.WaiGe)
	}
}

func TestScoreWuGeRange(t *testing.T) {
	// 全吉
	wg := WuGeResult{TianGe: 1, RenGe: 1, DiGe: 1, ZongGe: 1, WaiGe: 1}
	score := ScoreWuGe(wg)
	if score < 28 || score > 30 {
		t.Errorf("ScoreWuGe(all DaJi) = %.1f, want ~30", score)
	}

	// 全凶
	wg2 := WuGeResult{TianGe: 2, RenGe: 2, DiGe: 2, ZongGe: 2, WaiGe: 2}
	score2 := ScoreWuGe(wg2)
	if score2 > 5 {
		t.Errorf("ScoreWuGe(all DaXiong) = %.1f, want < 5", score2)
	}

	// 一般情况在 0-30 范围
	wg3 := WuGeResult{TianGe: 5, RenGe: 12, DiGe: 18, ZongGe: 22, WaiGe: 11}
	score3 := ScoreWuGe(wg3)
	if score3 < 0 || score3 > 30 {
		t.Errorf("ScoreWuGe = %.1f, out of [0, 30]", score3)
	}
}

func TestJiXiongToWuGeScore(t *testing.T) {
	tests := []struct {
		jx   data.JiXiong
		want float64
	}{
		{data.DaJi, 100},
		{data.ZhongJi, 90},
		{data.Ji, 85},
		{data.JiDuo, 65},
		{data.JiXiongBan, 50},
		{data.XiongDuo, 30},
		{data.DaXiong, 10},
		{data.JiXiong(99), 50}, // default
	}
	for _, tt := range tests {
		if got := jiXiongToWuGeScore(tt.jx); got != tt.want {
			t.Errorf("jiXiongToWuGeScore(%v) = %.0f, want %.0f", tt.jx, got, tt.want)
		}
	}
}
