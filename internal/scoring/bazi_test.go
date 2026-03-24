package scoring

import (
	"testing"

	"github.com/vogo/namer/internal/data"
)

func TestGanZhiString(t *testing.T) {
	gz := GanZhi{Gan: 0, Zhi: 0}
	if got := gz.String(); got != "甲子" {
		t.Errorf("GanZhi{0,0}.String() = %q, want %q", got, "甲子")
	}
	gz2 := GanZhi{Gan: 9, Zhi: 11}
	if got := gz2.String(); got != "癸亥" {
		t.Errorf("GanZhi{9,11}.String() = %q, want %q", got, "癸亥")
	}
}

func TestBaZiResultString(t *testing.T) {
	bz := BaZiResult{
		Year:  GanZhi{0, 0},
		Month: GanZhi{2, 2},
		Day:   GanZhi{4, 4},
		Hour:  GanZhi{6, 6},
	}
	got := bz.String()
	if got != "甲子年 丙寅月 戊辰日 庚午时" {
		t.Errorf("BaZiResult.String() = %q", got)
	}
}

func TestCalcYearGanZhi(t *testing.T) {
	tests := []struct {
		year    int
		wantGan int
		wantZhi int
	}{
		{2024, 0, 4},   // 甲辰
		{1984, 0, 0},   // 甲子
		{2000, 6, 4},   // 庚辰
		{1990, 6, 6},   // 庚午
		{2023, 9, 3},   // 癸卯
	}
	for _, tt := range tests {
		gz := calcYearGanZhi(tt.year)
		if gz.Gan != tt.wantGan || gz.Zhi != tt.wantZhi {
			t.Errorf("calcYearGanZhi(%d) = %v, want Gan=%d Zhi=%d", tt.year, gz, tt.wantGan, tt.wantZhi)
		}
	}
}

func TestCalcMonthGanZhi(t *testing.T) {
	tests := []struct {
		yearGan int
		month   int
		wantGan int
		wantZhi int
	}{
		// 甲年(0) 正月 → 丙寅(2, 2)
		{0, 1, 2, 2},
		// 甲年(0) 二月 → 丁卯(3, 3)
		{0, 2, 3, 3},
		// 乙年(1) 正月 → 戊寅(4, 2)
		{1, 1, 4, 2},
		// 丙年(2) 正月 → 庚寅(6, 2)
		{2, 1, 6, 2},
		// 丁年(3) 正月 → 壬寅(8, 2)
		{3, 1, 8, 2},
		// 戊年(4) 正月 → 甲寅(0, 2)
		{4, 1, 0, 2},
	}
	for _, tt := range tests {
		gz := calcMonthGanZhi(tt.yearGan, tt.month)
		if gz.Gan != tt.wantGan || gz.Zhi != tt.wantZhi {
			t.Errorf("calcMonthGanZhi(%d, %d) = %s, want Gan=%d Zhi=%d",
				tt.yearGan, tt.month, gz.String(), tt.wantGan, tt.wantZhi)
		}
	}
}

func TestCalcHourGanZhi(t *testing.T) {
	tests := []struct {
		dayGan  int
		hour    int
		wantGan int
		wantZhi int
	}{
		// 甲日(0) 子时(23-1, zhi=0) → 甲子(0,0)
		{0, 0, 0, 0},
		// 甲日(0) 寅时(3-5, zhi=2) → 丙寅(2,2)
		{0, 4, 2, 2},
		// 乙日(1) 子时 → 丙子(2,0)
		{1, 0, 2, 0},
	}
	for _, tt := range tests {
		gz := calcHourGanZhi(tt.dayGan, tt.hour)
		if gz.Gan != tt.wantGan || gz.Zhi != tt.wantZhi {
			t.Errorf("calcHourGanZhi(%d, %d) = %s, want Gan=%d Zhi=%d",
				tt.dayGan, tt.hour, gz.String(), tt.wantGan, tt.wantZhi)
		}
	}
}

func TestCalcBaZi(t *testing.T) {
	// 基本测试：确保不 panic，且各柱在有效范围
	bz := CalcBaZi(2024, 3, 15, 10)
	if bz.Year.Gan < 0 || bz.Year.Gan > 9 {
		t.Errorf("Year Gan out of range: %d", bz.Year.Gan)
	}
	if bz.Year.Zhi < 0 || bz.Year.Zhi > 11 {
		t.Errorf("Year Zhi out of range: %d", bz.Year.Zhi)
	}
	if bz.Month.Gan < 0 || bz.Month.Gan > 9 {
		t.Errorf("Month Gan out of range: %d", bz.Month.Gan)
	}
	if bz.Day.Gan < 0 || bz.Day.Gan > 9 {
		t.Errorf("Day Gan out of range: %d", bz.Day.Gan)
	}
	if bz.Hour.Gan < 0 || bz.Hour.Gan > 9 {
		t.Errorf("Hour Gan out of range: %d", bz.Hour.Gan)
	}

	// 2024 = 甲辰年
	if bz.Year.Gan != 0 || bz.Year.Zhi != 4 {
		t.Errorf("2024 year = %s, want 甲辰", bz.Year.String())
	}
}

func TestCalcBaZiJanFeb(t *testing.T) {
	// 测试1月2月（month <= 2 的特殊分支）
	bz := CalcBaZi(2024, 1, 15, 10)
	if bz.Year.Gan < 0 || bz.Year.Gan > 9 {
		t.Errorf("Jan: Year Gan out of range: %d", bz.Year.Gan)
	}

	bz2 := CalcBaZi(2024, 2, 15, 10)
	if bz2.Day.Gan < 0 || bz2.Day.Gan > 9 {
		t.Errorf("Feb: Day Gan out of range: %d", bz2.Day.Gan)
	}
}

func TestCalcXiYongShen(t *testing.T) {
	// 测试不同八字返回有效的五行
	years := []int{1990, 2000, 2010, 2024}
	months := []int{1, 6, 12}
	for _, y := range years {
		for _, m := range months {
			bz := CalcBaZi(y, m, 15, 12)
			xy := CalcXiYongShen(bz)
			if xy == data.WuXingUnknown {
				t.Errorf("CalcXiYongShen for %d-%d returned WuXingUnknown", y, m)
			}
			if xy < data.Jin || xy > data.Tu {
				t.Errorf("CalcXiYongShen for %d-%d returned invalid WuXing: %v", y, m, xy)
			}
		}
	}
}

func TestCalcDayGanZhi(t *testing.T) {
	// 验证不同日期的日柱在有效范围内
	dates := [][3]int{
		{2024, 1, 1}, {2024, 6, 15}, {2024, 12, 31},
		{2000, 1, 1}, {1990, 7, 20}, {2023, 3, 8},
	}
	for _, d := range dates {
		gz := calcDayGanZhi(d[0], d[1], d[2])
		if gz.Gan < 0 || gz.Gan > 9 || gz.Zhi < 0 || gz.Zhi > 11 {
			t.Errorf("calcDayGanZhi(%d,%d,%d) = Gan=%d Zhi=%d, out of range",
				d[0], d[1], d[2], gz.Gan, gz.Zhi)
		}
	}
}
