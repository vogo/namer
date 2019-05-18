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

最后, 执行命令:
```bash
namer -c config.json
```
命令会输出分数排名前10的名字列表。

## 注意
- 配置可以修改以后再次执行。
- 命令会生成缓存文件 `config.json.score`，请勿删除。