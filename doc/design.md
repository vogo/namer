# 姓名阴阳五行评分工具设计

## 概述

设计一个基于阴阳五行理论的姓名评分工具，满分100分。所有数据和算法内置于项目中，不依赖外部API。

## 输入

```json
{
  "last_name": "王",
  "first_name": "明轩",
  "year": 2024,
  "month": 3,
  "day": 15,
  "hour": 10,
  "minute": 30,
  "gender": 0
}
```

## 评分维度（总分100分）

| 维度                     | 分值 | 说明                         |
|--------------------------|------|------------------------------|
| 1. 五格数理吉凶          | 30分 | 天格、人格、地格、总格、外格  |
| 2. 三才五行配置          | 25分 | 天人地三才的五行生克关系      |
| 3. 名字五行与喜用神匹配  | 20分 | 名字五行是否补益八字喜用神    |
| 4. 名字内部五行生克      | 15分 | 姓与名各字之间的五行关系      |
| 5. 阴阳平衡             | 10分 | 姓名各字笔画的阴阳搭配        |

---

## 项目内置数据

所有数据直接以 Go 变量硬编码在源码中，不使用外部 JSON 配置文件，方便代码直接判断。

### data/strokes.go — 汉字康熙笔画数据

按笔画数分组，每组一个 `map[rune]bool`，同时提供查询函数。

```go
package data

// 1画的汉字
var Strokes1 = map[rune]bool{
    '一': true, '乙': true,
}

// 2画的汉字
var Strokes2 = map[rune]bool{
    '卜': true, '刀': true, '刁': true, '丁': true,
    '二': true, '力': true, '了': true, '人': true,
    '入': true, '七': true, '十': true, '几': true,
}

// ... 3画到24画，每个笔画数一个 map

// CharStroke 查询汉字的康熙笔画数，未找到返回 0
func CharStroke(c rune) int {
    if Strokes1[c] { return 1 }
    if Strokes2[c] { return 2 }
    // ... 依次查询各笔画 map
    return 0
}
```

来源：基于 `doc/相应笔画下常用汉字.md` 中已有的汉字笔画数据。笔画数按康熙字典计算（需按 `doc/计算五格数时要注意两个问题.md` 中的特殊偏旁规则处理）。

### data/wuxing_jin.go, wuxing_mu.go, wuxing_shui.go, wuxing_huo.go, wuxing_tu.go — 汉字五行属性

每个五行一个文件，各包含一个 `map[rune]bool`。

```go
package data

// wuxing_jin.go — 五行属金的汉字
var WuXingJin = map[rune]bool{
    '金': true, '钢': true, '铁': true, '银': true, '铜': true,
    '锋': true, '铭': true, '锐': true, '钊': true, '钧': true,
    '刚': true, '利': true, '列': true, '刘': true, '则': true,
    // ... 所有五行属金的字
}
```

```go
package data

// wuxing_mu.go — 五行属木的汉字
var WuXingMu = map[rune]bool{
    '木': true, '林': true, '森': true, '松': true, '柏': true,
    '杨': true, '柳': true, '梅': true, '桂': true, '栋': true,
    '芳': true, '花': true, '草': true, '荣': true, '莉': true,
    // ... 所有五行属木的字
}
```

```go
package data

// wuxing_shui.go — 五行属水的汉字
var WuXingShui = map[rune]bool{
    '水': true, '河': true, '海': true, '湖': true, '江': true,
    '波': true, '涛': true, '洋': true, '涵': true, '泽': true,
    // ... 所有五行属水的字
}
```

```go
package data

// wuxing_huo.go — 五行属火的汉字
var WuXingHuo = map[rune]bool{
    '火': true, '炎': true, '焱': true, '灿': true, '烁': true,
    '光': true, '明': true, '辉': true, '耀': true, '晖': true,
    // ... 所有五行属火的字
}
```

```go
package data

// wuxing_tu.go — 五行属土的汉字
var WuXingTu = map[rune]bool{
    '土': true, '山': true, '岩': true, '峰': true, '崇': true,
    '坤': true, '垣': true, '城': true, '坚': true, '培': true,
    // ... 所有五行属土的字
}
```

提供统一查询函数：

```go
package data

// WuXing 五行类型
type WuXing int

const (
    WuXingUnknown WuXing = iota
    Jin   // 金
    Mu    // 木
    Shui  // 水
    Huo   // 火
    Tu    // 土
)

// CharWuXing 查询汉字的五行属性
func CharWuXing(c rune) WuXing {
    if WuXingJin[c]  { return Jin }
    if WuXingMu[c]   { return Mu }
    if WuXingShui[c] { return Shui }
    if WuXingHuo[c]  { return Huo }
    if WuXingTu[c]   { return Tu }
    return WuXingUnknown
}
```

来源：基于 `doc/相应笔画下常用汉字.md` 中已有的汉字五行标注。判定优先级参照 `doc/汉字五行属性判定方法.md`：字义 > 字形 > 数理。

### data/sancai.go — 三才配置吉凶表

```go
package data

// JiXiong 吉凶等级
type JiXiong int

const (
    DaXiong    JiXiong = iota // 大凶
    XiongDuo                   // 凶多吉少 / 凶多于吉
    JiXiongBan                 // 吉凶参半
    JiDuo                      // 吉多于凶
    Ji                         // 吉
    ZhongJi                    // 中吉
    DaJi                       // 大吉
)

// SanCaiJiXiong 三才配置吉凶表，key 为 "天五行人五行地五行"，如 "木木木"
var SanCaiJiXiong = map[string]JiXiong{
    "木木木": DaJi,
    "木木火": DaJi,
    "木木土": DaJi,
    "木木金": XiongDuo,
    "木木水": JiDuo,
    "木火木": DaJi,
    "木火火": ZhongJi,
    "木火土": DaJi,
    // ... 全部125种配置
}
```

来源：`doc/三才五行配置吉凶表.md` 中全部125种三才配置。

### data/wuge_jixiong.go — 五格数理吉凶表

```go
package data

// WuGeJiXiong 1-81数理的吉凶，索引0不用，从1开始
// 数理超过81则取模（n % 80，0视为80）
var WuGeJiXiong = [82]JiXiong{
    0:  DaJi,     // 占位，不使用
    1:  DaJi,     // 宇宙起源
    2:  DaXiong,  // 混沌未定
    3:  DaJi,     // 万物成形
    4:  DaXiong,  // 日月无光
    5:  DaJi,     // 福禄长寿
    6:  DaJi,     // 天德地祥
    7:  Ji,       // 刚毅果断
    8:  Ji,       // 意志坚强
    // ... 9-81
}

// GetWuGeJiXiong 获取数理的吉凶，自动处理超过81的数理
func GetWuGeJiXiong(n int) JiXiong {
    if n <= 0 {
        return DaXiong
    }
    if n > 81 {
        n = n % 80
        if n == 0 {
            n = 80
        }
    }
    return WuGeJiXiong[n]
}
```

81个数理（1-81）的吉凶定义。数理超过81则取模（n % 80，0视为80）。

---

## 算法详细设计

### 步骤一：计算康熙笔画数

调用 `data.CharStroke()` 查询每个字的康熙笔画数。

```
示例：王(4) 明(8) 轩(10)
```

### 步骤二：计算五格

按 `doc/五格计算方法.md` 的规则：

**单姓双名（最常见）：**

```
天格 = 姓笔画 + 1
人格 = 姓笔画 + 名1笔画
地格 = 名1笔画 + 名2笔画
总格 = 姓笔画 + 名1笔画 + 名2笔画
外格 = 总格 - 人格 + 1
```

**单姓单名：**

```
天格 = 姓笔画 + 1
人格 = 姓笔画 + 名笔画
地格 = 名笔画 + 1
总格 = 姓笔画 + 名笔画
外格 = 2
```

**复姓双名：**

```
天格 = 姓1笔画 + 姓2笔画
人格 = 姓2笔画 + 名1笔画
地格 = 名1笔画 + 名2笔画
总格 = 全部笔画之和
外格 = 总格 - 人格
```

**复姓单名：**

```
天格 = 姓1笔画 + 姓2笔画
人格 = 姓2笔画 + 名笔画
地格 = 名笔画 + 1
总格 = 全部笔画之和
外格 = 总格 - 人格 + 1
```

### 步骤三：五格数理评分（30分）

每格数理调用 `data.GetWuGeJiXiong()` 查询吉凶，按以下权重计算：

| 格     | 权重 | 说明                     |
|--------|------|--------------------------|
| 人格   | 35%  | 主运，影响最大            |
| 地格   | 25%  | 前运，影响青少年           |
| 总格   | 20%  | 后运，影响中晚年           |
| 外格   | 15%  | 副运，人际关系             |
| 天格   | 5%   | 先天运，姓氏决定无法改变    |

每格的吉凶转换为分数：

| 吉凶     | 分数 |
|----------|------|
| 大吉     | 100  |
| 吉       | 85   |
| 半吉     | 65   |
| 凶       | 25   |
| 大凶     | 10   |

```
五格得分 = (人格分×0.35 + 地格分×0.25 + 总格分×0.20 + 外格分×0.15 + 天格分×0.05) / 100 × 30
```

### 步骤四：三才五行配置评分（25分）

1. 确定天格、人格、地格各自的五行（个位数: 1,2=木 3,4=火 5,6=土 7,8=金 9,0=水）
2. 查 `data.SanCaiJiXiong` 获取三才配置的吉凶
3. 转换为分数：

| 吉凶       | 分数 |
|------------|------|
| 大吉       | 100  |
| 中吉       | 80   |
| 吉         | 75   |
| 吉多于凶   | 60   |
| 吉凶参半   | 50   |
| 凶多吉少   | 30   |
| 凶多于吉   | 30   |
| 大凶       | 10   |

```
三才得分 = 三才分数 / 100 × 25
```

### 步骤五：名字五行与喜用神匹配评分（20分）

#### 5a. 计算八字

根据出生年月日时，通过天干地支算法计算四柱八字。

**年柱天干算法：**
```
天干索引 = (year - 4) % 10
地支索引 = (year - 4) % 12
```

天干序列：甲乙丙丁戊己庚辛壬癸
地支序列：子丑寅卯辰巳午未申酉戌亥

**月柱：** 根据年干和月份查表（五虎遁月法）。
**日柱：** 采用日期序号算法计算日干支。
**时柱：** 根据日干和时辰查表（五鼠遁时法）。

#### 5b. 确定喜用神

1. 统计八字中金木水火土各五行的个数和强度
2. 确定日主（日柱天干）的五行
3. 判断日主强弱
4. 日主强 → 喜用神为克泄日主的五行
5. 日主弱 → 喜用神为生扶日主的五行

**简化版喜用神算法：**

```go
func calcXiYongShen(bazi []GanZhi) WuXing {
    dayMaster := bazi[2].Gan  // 日柱天干
    dayElement := ganWuXing(dayMaster)

    // 统计各五行力量
    strength := countElementStrength(bazi)

    // 判断日主强弱
    if isDayMasterStrong(dayMaster, strength) {
        // 日主强，喜克泄
        return getWeakeningElement(dayElement)
    } else {
        // 日主弱，喜生扶
        return getStrengtheningElement(dayElement)
    }
}
```

#### 5c. 评分

调用 `data.CharWuXing()` 查询名字每个字的五行属性，检查与喜用神的关系：

| 关系                   | 得分率 |
|------------------------|--------|
| 名字五行 = 喜用神      | 100%   |
| 名字五行生喜用神       | 80%    |
| 喜用神生名字五行       | 60%    |
| 无直接关系             | 40%    |
| 名字五行克喜用神       | 20%    |
| 喜用神克名字五行       | 10%    |

多字取平均：

```
喜用神匹配得分 = avg(各字得分率) / 100 × 20
```

### 步骤六：名字内部五行生克评分（15分）

检查姓名相邻字之间的五行关系：

| 关系           | 分数 |
|----------------|------|
| 前字生后字     | 100  |
| 后字生前字     | 80   |
| 五行相同       | 70   |
| 无直接关系     | 50   |
| 前字克后字     | 20   |
| 后字克前字     | 30   |

对于三字姓名，有两组相邻关系（姓-名1, 名1-名2），取平均：

```
内部五行得分 = avg(各组分数) / 100 × 15
```

### 步骤七：阴阳平衡评分（10分）

根据每个字的康熙笔画数判断阴阳（奇数=阳，偶数=阴）：

**三字姓名：**

| 格局           | 分数 |
|----------------|------|
| 阳阴阳 / 阴阳阴 | 100  |
| 阴阳阳 / 阳阳阴 | 80   |
| 阳阴阴 / 阴阴阳 | 80   |
| 阳阳阳 / 阴阴阴 | 40   |

**两字姓名：**

| 格局       | 分数 |
|------------|------|
| 阳阴 / 阴阳 | 100  |
| 阳阳 / 阴阴 | 40   |

```
阴阳得分 = 阴阳分数 / 100 × 10
```

### 最终总分

```
总分 = 五格得分 + 三才得分 + 喜用神匹配得分 + 内部五行得分 + 阴阳得分
```

---

## 项目代码结构

```
namer/
├── doc/                          # 理论资料文档
│   ├── design.md                 # 本设计文档
│   ├── 五行基础.md
│   ├── 汉字五行属性判定方法.md
│   ├── 阴阳平衡.md
│   ├── 三才五行配置吉凶表.md
│   ├── 生辰八字与五行.md
│   ├── 五格计算方法.md
│   ├── 五格姓名影响.md
│   ├── 计算五格数时要注意两个问题.md
│   ├── 天地人格最佳搭配百家姓常用笔画数.md
│   └── 相应笔画下常用汉字.md
├── data/                         # 内置数据（Go源码硬编码）
│   ├── wuxing.go                 # WuXing 类型定义、JiXiong 类型定义、查询函数
│   ├── wuxing_jin.go             # WuXingJin map — 五行属金的汉字
│   ├── wuxing_mu.go              # WuXingMu map — 五行属木的汉字
│   ├── wuxing_shui.go            # WuXingShui map — 五行属水的汉字
│   ├── wuxing_huo.go             # WuXingHuo map — 五行属火的汉字
│   ├── wuxing_tu.go              # WuXingTu map — 五行属土的汉字
│   ├── strokes.go                # CharStroke() — 汉字康熙笔画查询（按笔画数分组的map）
│   ├── sancai.go                 # SanCaiJiXiong map — 三才配置吉凶表(125项)
│   └── wuge_jixiong.go           # WuGeJiXiong 数组 — 五格数理吉凶表(1-81)
├── bazi.go                       # 八字计算（年月日时四柱）
├── wuge.go                       # 五格计算（天人地总外）
├── sancai.go                     # 三才五行配置评分
├── wuxing.go                     # 五行生克关系 & 喜用神计算
├── yinyang.go                    # 阴阳平衡评分
├── score.go                      # 综合评分引擎（整合各维度）
├── namer.go                      # 主程序入口（已有，调整集成评分）
├── go.mod
└── README.md
```

## 核心接口设计

```go
// 评分结果
type ScoreResult struct {
    Total          float64          // 总分 (0-100)
    WuGeScore      float64          // 五格数理得分 (0-30)
    SanCaiScore    float64          // 三才配置得分 (0-25)
    XiYongScore    float64          // 喜用神匹配得分 (0-20)
    WuXingScore    float64          // 内部五行得分 (0-15)
    YinYangScore   float64          // 阴阳平衡得分 (0-10)
    Detail         ScoreDetail      // 详细计算过程
}

// 评分详情
type ScoreDetail struct {
    Strokes        []int            // 各字康熙笔画
    WuGe           WuGeResult       // 五格计算结果
    SanCai         SanCaiResult     // 三才配置结果
    BaZi           BaZiResult       // 八字结果
    XiYongShen     string           // 喜用神五行
    CharWuXing     []string         // 各字五行属性
    YinYangPattern string           // 阴阳格局
}

// 五格结果
type WuGeResult struct {
    TianGe  int    // 天格数
    RenGe   int    // 人格数
    DiGe    int    // 地格数
    ZongGe  int    // 总格数
    WaiGe   int    // 外格数
}

// 主评分函数
func CalcScore(lastName, firstName string, birth time.Time) ScoreResult
```

## 使用方式

```bash
# 评分模式
go run . -mode score -last 王 -first 明轩 -birth "2024-03-15 10:30"

# 输出示例:
# 姓名: 王明轩
# 总分: 87.5 / 100
# ┌──────────────┬──────┬──────┐
# │ 评分维度      │ 得分 │ 满分 │
# ├──────────────┼──────┼──────┤
# │ 五格数理      │ 25.2 │  30  │
# │ 三才配置      │ 25.0 │  25  │
# │ 喜用神匹配    │ 16.0 │  20  │
# │ 内部五行      │ 13.5 │  15  │
# │ 阴阳平衡      │  7.8 │  10  │
# └──────────────┴──────┴──────┘
#
# 详细信息:
# 康熙笔画: 王(4) 明(8) 轩(10)
# 五格: 天格5 人格12 地格18 总格22 外格11
# 三才: 土木金 → 凶多于吉
# 八字: 甲辰年 丁卯月 壬午日 乙巳时
# 喜用神: 金
# 字五行: 金 火 土
# 阴阳: 阴 阴 阴 → 纯阴
```
