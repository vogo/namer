# Namer 是中文取名工具，其根据你喜欢的字，自动组合出符合阴阳五行的名字。

## 如何使用

首先, 下载工具:
```bash
go get github.com/wongoo/namer
```

其次, 准备配置文件config.json:
```json
{
  "last_name": "王",
  "year": 2018,
  "month": 8,
  "day": 15,
  "hour": 11,
  "minute": 1,
  "gender": 0,
  "first_name_key_words": "可,计,学,习,书,意,义,复,开,程,进,望"
}
```

最后, 执行命令 `namer -c config.json`, 命令会输出分数排名前10的名字列表:
```bash
score: 93, names: [王程]

score: 92, names: [王计程]

score: 89, names: [王望进 王开]

score: 88, names: [王望程 王望开 王习程 王程意 王程义 王程进 王计开]

score: 87, names: [王程开 王程程 王开程 王计学]

score: 86, names: [王程习 王程望 王开望]

score: 85, names: [王望学 王望意 王望复 王望义 王习进 王开进]

score: 84, names: [王习开 王程复 王程学 王开意 王开义 王开开 王书进]

score: 83, names: [王义程 王意程]

score: 82, names: [王习意 王习义 王习复 王开习 王书望]
```

## 注意
- 配置可以修改以后再次执行;
- 可以`ctrl+c`中断执行,立即给出已经检测的结果;
- 命令会生成缓存文件 `config.json.score`，请勿删除;