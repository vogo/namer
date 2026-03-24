package scoring

import "testing"

func TestScoreYinYang(t *testing.T) {
	tests := []struct {
		name        string
		strokes     []int
		wantScore   float64
		wantPattern string
	}{
		{"单字", []int{5}, 5.0, "未知"},
		{"阳阴(交替)", []int{5, 8}, 10.0, "阳阴"},
		{"阴阳(交替)", []int{8, 5}, 10.0, "阴阳"},
		{"阳阳(纯阳)", []int{5, 7}, 4.0, "阳阳"},
		{"阴阴(纯阴)", []int{8, 6}, 4.0, "阴阴"},
		{"阳阴阳(交替)", []int{5, 8, 7}, 10.0, "阳阴阳"},
		{"阴阳阴(交替)", []int{8, 5, 6}, 10.0, "阴阳阴"},
		{"阳阳阳(纯阳)", []int{5, 7, 9}, 4.0, "阳阳阳"},
		{"阴阴阴(纯阴)", []int{8, 6, 4}, 4.0, "阴阴阴"},
		{"阳阳阴(有阴有阳)", []int{5, 7, 8}, 8.0, "阳阳阴"},
		{"阴阳阳(有阴有阳)", []int{8, 5, 7}, 8.0, "阴阳阳"},
		{"阴阴阳(有阴有阳)", []int{8, 6, 7}, 8.0, "阴阴阳"},
		{"阳阴阴(有阴有阳)", []int{5, 8, 6}, 8.0, "阳阴阴"},
	}
	for _, tt := range tests {
		score, pattern := ScoreYinYang(tt.strokes)
		if score != tt.wantScore {
			t.Errorf("%s: score = %.1f, want %.1f", tt.name, score, tt.wantScore)
		}
		if pattern != tt.wantPattern {
			t.Errorf("%s: pattern = %q, want %q", tt.name, pattern, tt.wantPattern)
		}
	}
}

func TestScoreYinYangFourChars(t *testing.T) {
	// 四字名 交替
	score, pattern := ScoreYinYang([]int{5, 8, 7, 6})
	if score != 10.0 {
		t.Errorf("alternating 4-char: score = %.1f, want 10.0", score)
	}
	if pattern != "阳阴阳阴" {
		t.Errorf("alternating 4-char: pattern = %q, want 阳阴阳阴", pattern)
	}

	// 四字名 纯阳
	score2, _ := ScoreYinYang([]int{1, 3, 5, 7})
	if score2 != 4.0 {
		t.Errorf("all yang 4-char: score = %.1f, want 4.0", score2)
	}

	// 四字名 混合
	score3, _ := ScoreYinYang([]int{1, 3, 4, 7})
	if score3 != 8.0 {
		t.Errorf("mixed 4-char: score = %.1f, want 8.0", score3)
	}
}
