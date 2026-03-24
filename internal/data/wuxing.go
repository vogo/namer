package data

// WuXing 五行类型
type WuXing int

const (
	WuXingUnknown WuXing = iota
	Jin                  // 金
	Mu                   // 木
	Shui                 // 水
	Huo                  // 火
	Tu                   // 土
)

// WuXingName 五行中文名
func (w WuXing) String() string {
	switch w {
	case Jin:
		return "金"
	case Mu:
		return "木"
	case Shui:
		return "水"
	case Huo:
		return "火"
	case Tu:
		return "土"
	default:
		return "未知"
	}
}

// CharWuXing 查询汉字的五行属性
func CharWuXing(c rune) WuXing {
	if WuXingJin[c] {
		return Jin
	}
	if WuXingMu[c] {
		return Mu
	}
	if WuXingShui[c] {
		return Shui
	}
	if WuXingHuo[c] {
		return Huo
	}
	if WuXingTu[c] {
		return Tu
	}
	return WuXingUnknown
}

// JiXiong 吉凶等级
type JiXiong int

const (
	DaXiong    JiXiong = iota // 大凶
	XiongDuo                  // 凶多吉少 / 凶多于吉
	JiXiongBan                // 吉凶参半
	JiDuo                     // 吉多于凶
	Ji                        // 吉
	ZhongJi                   // 中吉
	DaJi                      // 大吉
)

func (j JiXiong) String() string {
	switch j {
	case DaJi:
		return "大吉"
	case ZhongJi:
		return "中吉"
	case Ji:
		return "吉"
	case JiDuo:
		return "吉多于凶"
	case JiXiongBan:
		return "吉凶参半"
	case XiongDuo:
		return "凶多于吉"
	case DaXiong:
		return "大凶"
	default:
		return "未知"
	}
}

// NumToWuXing 五格数理的个位数转五行（五格剖象法）
func NumToWuXing(n int) WuXing {
	d := n % 10
	switch d {
	case 1, 2:
		return Mu
	case 3, 4:
		return Huo
	case 5, 6:
		return Tu
	case 7, 8:
		return Jin
	case 9, 0:
		return Shui
	}
	return WuXingUnknown
}

// Generates 五行相生：a 生 b
func Generates(a, b WuXing) bool {
	switch a {
	case Mu:
		return b == Huo
	case Huo:
		return b == Tu
	case Tu:
		return b == Jin
	case Jin:
		return b == Shui
	case Shui:
		return b == Mu
	}
	return false
}

// Overcomes 五行相克：a 克 b
func Overcomes(a, b WuXing) bool {
	switch a {
	case Mu:
		return b == Tu
	case Tu:
		return b == Shui
	case Shui:
		return b == Huo
	case Huo:
		return b == Jin
	case Jin:
		return b == Mu
	}
	return false
}

// GeneratingElement 返回生助 target 的五行
func GeneratingElement(target WuXing) WuXing {
	switch target {
	case Mu:
		return Shui
	case Huo:
		return Mu
	case Tu:
		return Huo
	case Jin:
		return Tu
	case Shui:
		return Jin
	}
	return WuXingUnknown
}

// WeakeningElement 返回克制 target 的五行
func WeakeningElement(target WuXing) WuXing {
	switch target {
	case Mu:
		return Jin
	case Huo:
		return Shui
	case Tu:
		return Mu
	case Jin:
		return Huo
	case Shui:
		return Tu
	}
	return WuXingUnknown
}
