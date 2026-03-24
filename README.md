# Namer

中文取名工具，根据你喜欢的字，自动组合出符合阴阳五行的名字。

> 也许你不相信阴阳五行，但工具产生出诸多好名字，会为你带来灵感，让你取一个满意的名字。

## 安装

从 [Releases](https://github.com/vogo/namer/releases) 页面下载对应平台的二进制文件（支持 macOS、Linux、Windows）。

或通过源码编译：

```bash
go install github.com/vogo/namer@latest
```

## 使用方法

### 交互式批量评分

```bash
namer
```

首次运行时会交互式引导创建配置文件（`~/.namer.json`），配置项包括姓氏、出生年月日时分、性别、名字备选字。

### 单个名字评分

```bash
namer <姓> <名>               # 使用已有配置中的生辰信息
namer <姓> <名> <生日>         # 指定生日，格式: 2024-03-15
```

示例：

```bash
namer 王 明轩
namer 王 明轩 2024-03-15
```

### 指定配置文件

```bash
namer -c <配置文件>
```

### Web 界面

```bash
namer -web                    # 随机端口，自动打开浏览器
namer -web -port 8080         # 指定端口
```

### 其他

```bash
namer -h                      # 显示帮助
namer -v                      # 显示版本
```

## 配置文件

默认路径：`~/.namer.json`，首次运行会交互式引导创建。

也可以手动创建 JSON 配置文件：

```json
{
  "last_name": "王",
  "year": 2024,
  "month": 3,
  "day": 15,
  "hour": 11,
  "minute": 1,
  "gender": 0,
  "min_candidate_score": 82,
  "first_name_key_words": "可,学,书,意,义,程,进,望"
}
```

| 字段 | 说明 |
|------|------|
| `last_name` | 姓 |
| `year` / `month` / `day` / `hour` / `minute` | 出生年月日时分 |
| `gender` | 0-男, 1-女 |
| `min_candidate_score` | 最小候选分数 |
| `first_name_key_words` | 名字备选字（逗号分隔） |

## 评分维度

总分 100 分：

| 维度 | 分值 | 说明 |
|------|------|------|
| 五格数理 | 30分 | 天格、人格、地格、总格、外格的数理吉凶 |
| 三才配置 | 25分 | 天人地三才的五行生克关系 |
| 喜用神 | 20分 | 名字五行是否补益八字喜用神 |
| 内部五行 | 15分 | 姓名各字之间的五行生克 |
| 阴阳平衡 | 10分 | 姓名各字笔画的阴阳搭配 |

## 输出示例

```
score: 93, names: [王程]
score: 92, names: [王计程]
score: 89, names: [王望进 王开]
score: 88, names: [王望程 王望开 王习程 王程意 王程义 王程进 王计开]
```

## 注意

- 配置可以修改后再次执行
- 可以 `Ctrl+C` 中断执行，立即给出已检测的结果
