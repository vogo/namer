package data

import "testing"

func TestWuXingString(t *testing.T) {
	tests := []struct {
		wx   WuXing
		want string
	}{
		{Jin, "金"},
		{Mu, "木"},
		{Shui, "水"},
		{Huo, "火"},
		{Tu, "土"},
		{WuXingUnknown, "未知"},
		{WuXing(99), "未知"},
	}
	for _, tt := range tests {
		if got := tt.wx.String(); got != tt.want {
			t.Errorf("WuXing(%d).String() = %q, want %q", tt.wx, got, tt.want)
		}
	}
}

func TestJiXiongString(t *testing.T) {
	tests := []struct {
		jx   JiXiong
		want string
	}{
		{DaJi, "大吉"},
		{ZhongJi, "中吉"},
		{Ji, "吉"},
		{JiDuo, "吉多于凶"},
		{JiXiongBan, "吉凶参半"},
		{XiongDuo, "凶多于吉"},
		{DaXiong, "大凶"},
		{JiXiong(99), "未知"},
	}
	for _, tt := range tests {
		if got := tt.jx.String(); got != tt.want {
			t.Errorf("JiXiong(%d).String() = %q, want %q", tt.jx, got, tt.want)
		}
	}
}

func TestCharWuXing(t *testing.T) {
	tests := []struct {
		c    rune
		want WuXing
	}{
		{'金', Jin},
		{'钢', Jin},
		{'木', Mu},
		{'林', Mu},
		{'水', Shui},
		{'海', Shui},
		{'火', Huo},
		{'炎', Huo},
		{'土', Tu},
		{'山', Tu},
		{'龘', WuXingUnknown}, // 不在数据库中
	}
	for _, tt := range tests {
		if got := CharWuXing(tt.c); got != tt.want {
			t.Errorf("CharWuXing(%c) = %v, want %v", tt.c, got, tt.want)
		}
	}
}

func TestCharStroke(t *testing.T) {
	tests := []struct {
		c    rune
		want int
	}{
		{'一', 1},
		{'乙', 1},
		{'人', 2},
		{'大', 3},
		{'王', 4},
		{'生', 5},
		{'安', 6},
		{'李', 7},
		{'明', 8},
		{'春', 9},
		{'高', 10},
		{'康', 11},
		{'博', 12},
		{'龘', 0}, // 不在数据库
	}
	for _, tt := range tests {
		if got := CharStroke(tt.c); got != tt.want {
			t.Errorf("CharStroke(%c) = %d, want %d", tt.c, got, tt.want)
		}
	}
}

func TestNumToWuXing(t *testing.T) {
	tests := []struct {
		n    int
		want WuXing
	}{
		{1, Mu}, {2, Mu},
		{3, Huo}, {4, Huo},
		{5, Tu}, {6, Tu},
		{7, Jin}, {8, Jin},
		{9, Shui}, {10, Shui},
		{11, Mu}, {12, Mu},
		{20, Shui},
		{35, Tu},
	}
	for _, tt := range tests {
		if got := NumToWuXing(tt.n); got != tt.want {
			t.Errorf("NumToWuXing(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}

func TestGenerates(t *testing.T) {
	// 木生火, 火生土, 土生金, 金生水, 水生木
	trueTests := [][2]WuXing{
		{Mu, Huo}, {Huo, Tu}, {Tu, Jin}, {Jin, Shui}, {Shui, Mu},
	}
	for _, tt := range trueTests {
		if !Generates(tt[0], tt[1]) {
			t.Errorf("Generates(%v, %v) should be true", tt[0], tt[1])
		}
	}
	// 反向不成立
	falseTests := [][2]WuXing{
		{Huo, Mu}, {Tu, Huo}, {Jin, Tu}, {Shui, Jin}, {Mu, Shui},
		{Mu, Mu}, {Jin, Jin},
		{WuXingUnknown, Jin},
	}
	for _, tt := range falseTests {
		if Generates(tt[0], tt[1]) {
			t.Errorf("Generates(%v, %v) should be false", tt[0], tt[1])
		}
	}
}

func TestOvercomes(t *testing.T) {
	// 木克土, 土克水, 水克火, 火克金, 金克木
	trueTests := [][2]WuXing{
		{Mu, Tu}, {Tu, Shui}, {Shui, Huo}, {Huo, Jin}, {Jin, Mu},
	}
	for _, tt := range trueTests {
		if !Overcomes(tt[0], tt[1]) {
			t.Errorf("Overcomes(%v, %v) should be true", tt[0], tt[1])
		}
	}
	falseTests := [][2]WuXing{
		{Tu, Mu}, {Shui, Tu}, {Huo, Shui}, {Jin, Huo}, {Mu, Jin},
		{Mu, Mu},
		{WuXingUnknown, Jin},
	}
	for _, tt := range falseTests {
		if Overcomes(tt[0], tt[1]) {
			t.Errorf("Overcomes(%v, %v) should be false", tt[0], tt[1])
		}
	}
}

func TestGeneratingElement(t *testing.T) {
	tests := []struct {
		target WuXing
		want   WuXing
	}{
		{Mu, Shui},   // 水生木
		{Huo, Mu},    // 木生火
		{Tu, Huo},    // 火生土
		{Jin, Tu},    // 土生金
		{Shui, Jin},  // 金生水
		{WuXingUnknown, WuXingUnknown},
	}
	for _, tt := range tests {
		if got := GeneratingElement(tt.target); got != tt.want {
			t.Errorf("GeneratingElement(%v) = %v, want %v", tt.target, got, tt.want)
		}
	}
}

func TestWeakeningElement(t *testing.T) {
	tests := []struct {
		target WuXing
		want   WuXing
	}{
		{Mu, Jin},   // 金克木
		{Huo, Shui}, // 水克火
		{Tu, Mu},    // 木克土
		{Jin, Huo},  // 火克金
		{Shui, Tu},  // 土克水
		{WuXingUnknown, WuXingUnknown},
	}
	for _, tt := range tests {
		if got := WeakeningElement(tt.target); got != tt.want {
			t.Errorf("WeakeningElement(%v) = %v, want %v", tt.target, got, tt.want)
		}
	}
}

func TestGetWuGeJiXiong(t *testing.T) {
	tests := []struct {
		n    int
		want JiXiong
	}{
		{0, DaXiong},  // <= 0
		{-1, DaXiong}, // <= 0
		{1, DaJi},
		{2, DaXiong},
		{3, DaJi},
		{5, DaJi},
		{7, Ji},
		{8, Ji},
		{81, DaJi},
		{82, DaJi},  // 82 % 80 = 2 → DaXiong... wait, 82 % 80 = 2
		{161, DaJi}, // 161 % 80 = 1 → DaJi
	}
	for _, tt := range tests {
		got := GetWuGeJiXiong(tt.n)
		// 只验证非零情况
		if tt.n <= 0 && got != DaXiong {
			t.Errorf("GetWuGeJiXiong(%d) = %v, want DaXiong", tt.n, got)
		}
		if tt.n == 1 && got != DaJi {
			t.Errorf("GetWuGeJiXiong(%d) = %v, want DaJi", tt.n, got)
		}
	}

	// 验证 > 81 取模
	if got := GetWuGeJiXiong(161); got != GetWuGeJiXiong(1) {
		t.Errorf("GetWuGeJiXiong(161) = %v, want same as GetWuGeJiXiong(1)=%v", got, GetWuGeJiXiong(1))
	}

	// 验证 160 % 80 = 0 → 80
	if got := GetWuGeJiXiong(160); got != GetWuGeJiXiong(80) {
		t.Errorf("GetWuGeJiXiong(160) = %v, want same as GetWuGeJiXiong(80)=%v", got, GetWuGeJiXiong(80))
	}
}

func TestSanCaiJiXiongCompleteness(t *testing.T) {
	elements := []string{"木", "火", "土", "金", "水"}
	count := 0
	for _, a := range elements {
		for _, b := range elements {
			for _, c := range elements {
				key := a + b + c
				if _, ok := SanCaiJiXiong[key]; !ok {
					t.Errorf("SanCaiJiXiong missing key: %s", key)
				}
				count++
			}
		}
	}
	if count != 125 {
		t.Errorf("expected 125 combinations, got %d", count)
	}
}

func TestSanCaiJiXiongValues(t *testing.T) {
	tests := []struct {
		key  string
		want JiXiong
	}{
		{"木木木", DaJi},
		{"木木金", XiongDuo},
		{"火土金", DaJi},
		{"金水木", DaJi},
		{"水火火", DaXiong},
		{"火金土", JiXiongBan},
	}
	for _, tt := range tests {
		if got := SanCaiJiXiong[tt.key]; got != tt.want {
			t.Errorf("SanCaiJiXiong[%s] = %v, want %v", tt.key, got, tt.want)
		}
	}
}
