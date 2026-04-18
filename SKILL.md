---
name: namer
description: "中文取名评分工具，根据阴阳五行原理自动组合并评分中文名字。用于安装 namer 命令行工具、对单个名字评分、批量组合评分。当用户需要中文取名、姓名评分、五行分析时使用。"
metadata:
  version: 0.1.0
  repository: https://github.com/vogo/namer
---

# 安装

从 [Releases](https://github.com/vogo/namer/releases) 页面下载对应平台的二进制文件（支持 macOS、Linux、Windows）。

或通过 Go 源码编译安装：

```bash
go install github.com/vogo/namer@latest
```

安装完成后，确保 `$GOPATH/bin` 在 `PATH` 中。验证安装：

```bash
namer -h
```

如果提示 `command not found`，检查 `$GOPATH/bin` 是否在 `PATH` 中：`export PATH=$PATH:$(go env GOPATH)/bin`。

# 配置文件

默认路径：`~/.namer.json`。首次运行 `namer`（无参数）时会交互式引导创建。也可手动创建：

```json
{
  "xing": "王",
  "year": 2024, "month": 3, "day": 15, "hour": 10, "minute": 30,
  "gender": 0,
  "min_candidate_score": 80,
  "ming_keywords": "明,轩,浩,然"
}
```

配置文件中的出生信息会作为命令行参数的默认值，命令行参数优先级更高。

# 使用示例

工具提供两大核心功能：**名字评分**（对指定名字打分）和**批量生成建议名字**（根据备选字自动组合并评分排序）。评分总分 100 分，涵盖五格数理（30分）、三才配置（25分）、喜用神匹配（20分）、内部五行（15分）、阴阳平衡（10分）。

## 一、名字评分

对已有名字进行五行评分分析。

```bash
namer -xing 王 -ming 明轩 -year 2024 -month 3 -day 15 -hour 10 -minute 30 -gender 1
```

如果配置文件（`~/.namer.json`）中已有生辰信息，可省略出生参数：

```bash
namer -xing 王 -ming 明轩
```

**参数说明：**

| 参数 | 说明 | 示例 |
|------|------|------|
| `-xing` | 姓氏（支持单姓和复姓） | `王`、`欧阳` |
| `-ming` | 名字（单字或双字） | `轩`、`明轩` |
| `-year -month -day -hour -minute` | 出生时间（省略时从配置读取） | `2024 3 15 10 30` |
| `-gender` | 性别（1=男, 2=女） | `1` |

输出包含总分、各维度得分、康熙笔画、五格、三才配置、八字、喜用神等详细分析。

## 二、批量生成建议名字

根据备选字自动组合（单字名 + 双字名），评分排序后输出高分名字。

```bash
namer -xing 王 -keywords 明,轩,浩,然 -year 2024 -month 3 -day 15 -hour 10 -minute 30 -gender 1 -score 70
```

| 参数 | 说明 | 示例 |
|------|------|------|
| `-keywords` | 名字备选字（逗号分隔） | `明,轩,浩,然` |
| `-score` | 最低分数过滤（默认无过滤） | `70` |

输出示例：

```
========== Top 10 ==========
 1. 王然浩  81.8 分
 2. 王浩然  80.4 分
 3. 王然  74.0 分
 4. 王浩轩  72.8 分
 5. 王然然  70.3 分
```
