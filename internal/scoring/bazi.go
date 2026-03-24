package scoring

import "github.com/vogo/namer/internal/data"

// 天干
var tianGan = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

// 地支
var diZhi = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

// GanZhi 干支
type GanZhi struct {
	Gan int // 天干索引 0-9
	Zhi int // 地支索引 0-11
}

func (gz GanZhi) String() string {
	return tianGan[gz.Gan] + diZhi[gz.Zhi]
}

// BaZiResult 八字结果
type BaZiResult struct {
	Year  GanZhi
	Month GanZhi
	Day   GanZhi
	Hour  GanZhi
}

func (b BaZiResult) String() string {
	return b.Year.String() + "年 " + b.Month.String() + "月 " + b.Day.String() + "日 " + b.Hour.String() + "时"
}

// 天干五行
var ganWuXing = []data.WuXing{
	data.Mu, data.Mu, // 甲乙
	data.Huo, data.Huo, // 丙丁
	data.Tu, data.Tu, // 戊己
	data.Jin, data.Jin, // 庚辛
	data.Shui, data.Shui, // 壬癸
}

// 地支五行
var zhiWuXing = []data.WuXing{
	data.Shui, // 子
	data.Tu,   // 丑
	data.Mu,   // 寅
	data.Mu,   // 卯
	data.Tu,   // 辰
	data.Huo,  // 巳
	data.Huo,  // 午
	data.Tu,   // 未
	data.Jin,  // 申
	data.Jin,  // 酉
	data.Tu,   // 戌
	data.Shui, // 亥
}

// 地支藏干：本气
var zhiCangGan = [][]int{
	{9},       // 子: 癸
	{5, 9, 7}, // 丑: 己癸辛
	{0, 2, 4}, // 寅: 甲丙戊
	{1},       // 卯: 乙
	{4, 1, 9}, // 辰: 戊乙癸
	{2, 4, 6}, // 巳: 丙戊庚
	{3, 5},    // 午: 丁己
	{5, 3, 1}, // 未: 己丁乙
	{6, 8, 4}, // 申: 庚壬戊
	{7},       // 酉: 辛
	{4, 7, 3}, // 戌: 戊辛丁
	{8, 0},    // 亥: 壬甲
}

// CalcBaZi 计算八字四柱
func CalcBaZi(year, month, day, hour int) BaZiResult {
	yearGZ := calcYearGanZhi(year)
	monthGZ := calcMonthGanZhi(yearGZ.Gan, month)
	dayGZ := calcDayGanZhi(year, month, day)
	hourGZ := calcHourGanZhi(dayGZ.Gan, hour)
	return BaZiResult{Year: yearGZ, Month: monthGZ, Day: dayGZ, Hour: hourGZ}
}

// 年柱
func calcYearGanZhi(year int) GanZhi {
	return GanZhi{
		Gan: (year - 4) % 10,
		Zhi: (year - 4) % 12,
	}
}

// 月柱 - 五虎遁月法
func calcMonthGanZhi(yearGan, month int) GanZhi {
	// 月地支固定：正月=寅(2), 二月=卯(3), ...
	monthZhi := (month + 1) % 12 // 正月=寅=2

	// 五虎遁月：根据年干确定正月天干
	// 甲己年起丙寅, 乙庚年起戊寅, 丙辛年起庚寅, 丁壬年起壬寅, 戊癸年起甲寅
	var startGan int
	switch yearGan % 5 {
	case 0: // 甲/己
		startGan = 2 // 丙
	case 1: // 乙/庚
		startGan = 4 // 戊
	case 2: // 丙/辛
		startGan = 6 // 庚
	case 3: // 丁/壬
		startGan = 8 // 壬
	case 4: // 戊/癸
		startGan = 0 // 甲
	}
	monthGan := (startGan + month - 1) % 10

	return GanZhi{Gan: monthGan, Zhi: monthZhi}
}

// 日柱 - 基于高氏日柱公式的简化算法
func calcDayGanZhi(year, month, day int) GanZhi {
	// 使用儒略日数计算
	if month <= 2 {
		year--
		month += 12
	}
	a := year / 100
	b := 2 - a + a/4
	jdn := int(365.25*float64(year+4716)) + int(30.6001*float64(month+1)) + day + b - 1524

	// 天干：(JDN - 1) % 10 映射到甲(0)
	// 已知 JDN 2451911 (2001-01-01) 的天干为 戊(4)
	gan := (jdn + 9) % 10
	if gan < 0 {
		gan += 10
	}

	// 地支：(JDN - 1) % 12
	// 已知 JDN 2451911 (2001-01-01) 的地支为 寅(2)
	zhi := (jdn + 1) % 12
	if zhi < 0 {
		zhi += 12
	}

	return GanZhi{Gan: gan, Zhi: zhi}
}

// 时柱 - 五鼠遁时法
func calcHourGanZhi(dayGan, hour int) GanZhi {
	// 时辰地支：23-1=子(0), 1-3=丑(1), ..., 21-23=亥(11)
	hourZhi := ((hour + 1) / 2) % 12

	// 五鼠遁时：根据日干确定子时天干
	// 甲己日起甲子, 乙庚日起丙子, 丙辛日起戊子, 丁壬日起庚子, 戊癸日起壬子
	var startGan int
	switch dayGan % 5 {
	case 0: // 甲/己
		startGan = 0 // 甲
	case 1: // 乙/庚
		startGan = 2 // 丙
	case 2: // 丙/辛
		startGan = 4 // 戊
	case 3: // 丁/壬
		startGan = 6 // 庚
	case 4: // 戊/癸
		startGan = 8 // 壬
	}
	hourGan := (startGan + hourZhi) % 10

	return GanZhi{Gan: hourGan, Zhi: hourZhi}
}

// CalcXiYongShen 计算喜用神
func CalcXiYongShen(bz BaZiResult) data.WuXing {
	dayMaster := ganWuXing[bz.Day.Gan]

	// 统计各五行力量
	strength := map[data.WuXing]int{
		data.Jin: 0, data.Mu: 0, data.Shui: 0, data.Huo: 0, data.Tu: 0,
	}

	// 四柱天干
	pillars := []GanZhi{bz.Year, bz.Month, bz.Day, bz.Hour}
	for _, p := range pillars {
		strength[ganWuXing[p.Gan]] += 10
	}

	// 四柱地支（本气权重较高）
	for _, p := range pillars {
		cangGans := zhiCangGan[p.Zhi]
		for i, g := range cangGans {
			w := ganWuXing[g]
			if i == 0 {
				strength[w] += 10 // 本气
			} else {
				strength[w] += 3 // 中气/余气
			}
		}
	}

	// 月令加权（月支本气为令）
	monthElement := zhiWuXing[bz.Month.Zhi]
	strength[monthElement] += 20

	// 计算日主同类（生我、同我）和异类（克我、我克、泄我）力量
	sameForce := strength[dayMaster] + strength[data.GeneratingElement(dayMaster)]
	diffForce := 0
	for wx, s := range strength {
		if wx != dayMaster && wx != data.GeneratingElement(dayMaster) {
			diffForce += s
		}
	}

	// 日主强弱判断
	if sameForce > diffForce {
		// 日主强，需要克泄耗 → 喜用神为克制或泄耗日主的五行
		// 优先选克，其次选泄
		overcoming := data.WeakeningElement(dayMaster)
		return overcoming
	}

	// 日主弱，需要生扶 → 喜用神为生助日主的五行
	return data.GeneratingElement(dayMaster)
}
